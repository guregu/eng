// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"image/color"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kardianos/osext"
)

var screen *ebiten.Image

var (
	responder Responder
	Files     *Loader
	bgColor   color.Color
)

func Open(title string, width, height int, fullscreen bool, r Responder) {
	responder = r
	// Time = NewClock()
	Files = NewLoader()
	run(title, width, height, fullscreen)
}

func SetBg(c uint32) {
	r := uint8((c >> 16) & 0xFF)
	g := uint8((c >> 8) & 0xFF)
	b := uint8(c & 0xFF)
	bgColor = color.RGBA{r, g, b, 255}
}

func Width() float64 {
	w, _ := ebiten.WindowSize()
	return float64(w)
}

func Height() float64 {
	_, h := ebiten.WindowSize()
	return float64(h)
}

func Exit() {
	os.Exit(0)
}

type ebitenGame struct {
	r Responder
	w int
	h int

	cursorX int
	cursorY int

	mousePressed bool
	prev         int64

	keys map[ebiten.Key]struct{}
}

func (g *ebitenGame) Update() error {
	cx, cy := ebiten.CursorPosition()
	if g.cursorX != cx || g.cursorY != cy {
		g.r.Mouse(cx, cy, MOVE)
	}

	mouse := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if !g.mousePressed && mouse {
		g.r.Mouse(cx, cy, PRESS)
	} else if g.mousePressed && !mouse {
		g.r.Mouse(cx, cy, RELEASE)
	}
	g.mousePressed = mouse

	var mod Modifier
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		mod |= SHIFT
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		mod |= CONTROL
	}
	if ebiten.IsKeyPressed(ebiten.KeyAlt) {
		mod |= ALT
	}
	for i := ebiten.Key(0); i <= ebiten.KeyMax; i++ {
		k := ebiten.Key(i)
		_, prev := g.keys[k]
		next := ebiten.IsKeyPressed(k)
		if next {
			g.keys[k] = struct{}{}
		} else {
			delete(g.keys, k)
		}

		if !prev && next {
			g.r.Key(Key(k), mod, PRESS)
		} else if prev && !next {
			g.r.Key(Key(k), mod, RELEASE)
		}
	}

	// TODO
	// now := time.Now().UnixNano()
	// dt := float64(((now - g.prev) / 1000000)) * 0.001
	// g.prev = now
	// // dt := 1.0 / ebiten.CurrentTPS()

	// g.r.Update(dt)
	g.r.Update(1.0 / float64(ebiten.MaxTPS()))
	return nil
}

func (g *ebitenGame) Draw(scr *ebiten.Image) {
	screen = scr
	scr.Fill(bgColor)
	g.r.Render()
}

func (g *ebitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.w, g.h
}

func run(title string, width, height int, fullscreen bool) {
	g := &ebitenGame{
		r:    responder,
		w:    width,
		h:    height,
		prev: time.Now().UnixNano(),
		keys: make(map[ebiten.Key]struct{}),
	}
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(width, height)
	ebiten.SetFullscreen(fullscreen)
	ebiten.SetWindowDecorated(!fullscreen)
	// ebiten.SetVsyncEnabled(true)
	// ebiten.SetWindowDecorated(!fullscreen)
	// ebiten.SetMaxTPS(ebiten.UncappedTPS)
	// ebiten.SetScreenClearedEveryFrame(true)

	// window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
	// 	responder.Scroll(float32(yoff))
	// })

	// window.SetCharCallback(func(window *glfw.Window, char rune) {
	// 	responder.Type(char)
	// })

	responder.Preload()
	// Files.Load(func() {})
	responder.Setup()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}

	responder.Close()
	if OnClose != nil {
		OnClose()
	}
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var OnClose func()

func Minimize() {
	ebiten.MinimizeWindow()
}

func AppDir() string {
	if runtime.GOOS == "js" {
		return ""
	}
	dir, err := osext.ExecutableFolder()
	fatalErr(err)
	return dir
}
