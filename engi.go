// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"image/color"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// import webgl "engo.io/gl"

var screen *ebiten.Image

var (
	responder Responder
	Time      *Clock
	Files     *Loader
	bgColor   color.Color
	// gl        *webgl.Context
)

func Open(title string, width, height int, fullscreen bool, r Responder) {
	responder = r
	Time = NewClock()
	Files = NewLoader()
	run(title, width, height, fullscreen)
}

func SetBg(c uint32) {
	// bgColor = c
	// TODO
	r := uint8((c >> 16) & 0xFF)
	g := uint8((c >> 8) & 0xFF)
	b := uint8(c & 0xFF)
	bgColor = color.RGBA{r, g, b, 255}
}

func Width() float32 {
	w, _ := ebiten.WindowSize()
	return float32(w)
}

func Height() float32 {
	_, h := ebiten.WindowSize()
	return float32(h)
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
		g.r.Mouse(float32(cx), float32(cy), MOVE)
	}

	mouse := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if !g.mousePressed && mouse {
		g.r.Mouse(float32(cx), float32(cy), PRESS)
	} else if g.mousePressed && !mouse {
		g.r.Mouse(float32(cx), float32(cy), RELEASE)
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
	// dt := 1.0 / ebiten.CurrentTPS()

	g.r.Update(1.0 / float32(ebiten.MaxTPS()))
	return nil
}

// var

func (g *ebitenGame) Draw(scr *ebiten.Image) {
	screen = scr
	scr.Fill(bgColor)
	g.r.Render() // TODO
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

	// window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
	// 	responder.Scroll(float32(yoff))
	// })

	// window.SetCharCallback(func(window *glfw.Window, char rune) {
	// 	responder.Type(char)
	// })

	// setupAudio()

	responder.Preload()
	Files.Load(func() {})
	responder.Setup()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}

	responder.Close()
	if OnClose != nil {
		OnClose()
	}
}
