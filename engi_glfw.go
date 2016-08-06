// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !netgo,!android

package engi

import (
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	webgl "engo.io/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/kardianos/osext"
)

var window *glfw.Window

func init() {
	runtime.LockOSThread()
}

// fatalErr calls log.Fatal with the given error if it is non-nil.
func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run(title string, width, height int, fullscreen bool) {
	defer runtime.UnlockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())
	fatalErr(glfw.Init())

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	// ideally we want a "windowless full screen"
	// but using normal full screen and disabling minimizing on OS X (which is broken)
	// is the best we can do for now...

	if fullscreen {
		// if runtime.GOOS == "darwin" {
		// 	glfw.WindowHint(glfw.AutoIconify, glfw.False)
		// }
		glfw.WindowHint(glfw.Decorated, glfw.False)
		//width, height = mode.Width, mode.Height
	} else {
		monitor = nil
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

	var err error
	window, err = glfw.CreateWindow(width, height, title, monitor, nil)
	fatalErr(err)
	window.MakeContextCurrent()

	if !fullscreen {
		window.SetPos((mode.Width-width)/2, (mode.Height-height)/2)
	}

	width, height = window.GetFramebufferSize()

	glfw.SwapInterval(1)

	gl = webgl.NewContext()

	gl.Viewport(0, 0, width, height)

	window.SetFramebufferSizeCallback(func(window *glfw.Window, w, h int) {
		oldWidth, oldHeight := width, height
		width, height = window.GetFramebufferSize()
		if runtime.GOOS == "darwin" && oldWidth/2 == w && oldHeight/2 == h {
			// I have no idea why, but this fixes fullscreen on OS X w/ retina
			gl.Viewport(0, h, w, h)
		} else {
			gl.Viewport(0, 0, w, h)
		}
		responder.Resize(w, h)
	})

	window.SetCursorPosCallback(func(window *glfw.Window, x, y float64) {
		responder.Mouse(float32(x), float32(y), MOVE)
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPos()
		if a == glfw.Press {
			responder.Mouse(float32(x), float32(y), PRESS)
		} else {
			responder.Mouse(float32(x), float32(y), RELEASE)
		}
	})

	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		responder.Scroll(float32(yoff))
	})

	window.SetKeyCallback(func(window *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		if a == glfw.Press {
			responder.Key(Key(k), Modifier(m), PRESS)
		} else if a == glfw.Release {
			responder.Key(Key(k), Modifier(m), RELEASE)
		}
	})

	window.SetCharCallback(func(window *glfw.Window, char rune) {
		responder.Type(char)
	})

	setupAudio()

	responder.Preload()
	Files.Load(func() {})
	responder.Setup()

	shouldClose := window.ShouldClose()
	for !shouldClose {
		responder.Update(Time.Delta())
		gl.Clear(gl.COLOR_BUFFER_BIT)
		responder.Render()
		window.SwapBuffers()
		glfw.PollEvents()
		Time.Tick()

		shouldClose = window.ShouldClose()
	}
	window.Destroy()
	glfw.Terminate()
	responder.Close()
}

func width() float32 {
	width, _ := window.GetSize()
	return float32(width)
}

func height() float32 {
	_, height := window.GetSize()
	return float32(height)
}

func exit() {
	cleanupAudio()
	window.SetShouldClose(true)
}

func init() {
	Dash = Key(glfw.KeyMinus)
	Apostrophe = Key(glfw.KeyApostrophe)
	Semicolon = Key(glfw.KeySemicolon)
	Equals = Key(glfw.KeyEqual)
	Comma = Key(glfw.KeyComma)
	Period = Key(glfw.KeyPeriod)
	Slash = Key(glfw.KeySlash)
	Backslash = Key(glfw.KeyBackslash)
	Backspace = Key(glfw.KeyBackspace)
	Tab = Key(glfw.KeyTab)
	CapsLock = Key(glfw.KeyCapsLock)
	Space = Key(glfw.KeySpace)
	Enter = Key(glfw.KeyEnter)
	Escape = Key(glfw.KeyEscape)
	Insert = Key(glfw.KeyInsert)
	PrintScreen = Key(glfw.KeyPrintScreen)
	Delete = Key(glfw.KeyDelete)
	PageUp = Key(glfw.KeyPageUp)
	PageDown = Key(glfw.KeyPageDown)
	Home = Key(glfw.KeyHome)
	End = Key(glfw.KeyEnd)
	Pause = Key(glfw.KeyPause)
	ScrollLock = Key(glfw.KeyScrollLock)
	ArrowLeft = Key(glfw.KeyLeft)
	ArrowRight = Key(glfw.KeyRight)
	ArrowDown = Key(glfw.KeyDown)
	ArrowUp = Key(glfw.KeyUp)
	LeftBracket = Key(glfw.KeyLeftBracket)
	LeftShift = Key(glfw.KeyLeftShift)
	LeftControl = Key(glfw.KeyLeftControl)
	LeftSuper = Key(glfw.KeyLeftSuper)
	LeftAlt = Key(glfw.KeyLeftAlt)
	RightBracket = Key(glfw.KeyRightBracket)
	RightShift = Key(glfw.KeyRightShift)
	RightControl = Key(glfw.KeyRightControl)
	RightSuper = Key(glfw.KeyRightSuper)
	RightAlt = Key(glfw.KeyRightAlt)
	Zero = Key(glfw.Key0)
	One = Key(glfw.Key1)
	Two = Key(glfw.Key2)
	Three = Key(glfw.Key3)
	Four = Key(glfw.Key4)
	Five = Key(glfw.Key5)
	Six = Key(glfw.Key6)
	Seven = Key(glfw.Key7)
	Eight = Key(glfw.Key8)
	Nine = Key(glfw.Key9)
	F1 = Key(glfw.KeyF1)
	F2 = Key(glfw.KeyF2)
	F3 = Key(glfw.KeyF3)
	F4 = Key(glfw.KeyF4)
	F5 = Key(glfw.KeyF5)
	F6 = Key(glfw.KeyF6)
	F7 = Key(glfw.KeyF7)
	F8 = Key(glfw.KeyF8)
	F9 = Key(glfw.KeyF9)
	F10 = Key(glfw.KeyF10)
	F11 = Key(glfw.KeyF11)
	F12 = Key(glfw.KeyF12)
	A = Key(glfw.KeyA)
	B = Key(glfw.KeyB)
	C = Key(glfw.KeyC)
	D = Key(glfw.KeyD)
	E = Key(glfw.KeyE)
	F = Key(glfw.KeyF)
	G = Key(glfw.KeyG)
	H = Key(glfw.KeyH)
	I = Key(glfw.KeyI)
	J = Key(glfw.KeyJ)
	K = Key(glfw.KeyK)
	L = Key(glfw.KeyL)
	M = Key(glfw.KeyM)
	N = Key(glfw.KeyN)
	O = Key(glfw.KeyO)
	P = Key(glfw.KeyP)
	Q = Key(glfw.KeyQ)
	R = Key(glfw.KeyR)
	S = Key(glfw.KeyS)
	T = Key(glfw.KeyT)
	U = Key(glfw.KeyU)
	V = Key(glfw.KeyV)
	W = Key(glfw.KeyW)
	X = Key(glfw.KeyX)
	Y = Key(glfw.KeyY)
	Z = Key(glfw.KeyZ)
	NumLock = Key(glfw.KeyNumLock)
	NumMultiply = Key(glfw.KeyKPMultiply)
	NumDivide = Key(glfw.KeyKPDivide)
	NumAdd = Key(glfw.KeyKPAdd)
	NumSubtract = Key(glfw.KeyKPSubtract)
	NumZero = Key(glfw.KeyKP0)
	NumOne = Key(glfw.KeyKP1)
	NumTwo = Key(glfw.KeyKP2)
	NumThree = Key(glfw.KeyKP3)
	NumFour = Key(glfw.KeyKP4)
	NumFive = Key(glfw.KeyKP5)
	NumSix = Key(glfw.KeyKP6)
	NumSeven = Key(glfw.KeyKP7)
	NumEight = Key(glfw.KeyKP8)
	NumNine = Key(glfw.KeyKP9)
	NumDecimal = Key(glfw.KeyKPDecimal)
	NumEnter = Key(glfw.KeyKPEnter)
}

func NewImageObject(img *image.NRGBA) *ImageObject {
	return &ImageObject{img}
}

type ImageObject struct {
	data *image.NRGBA
}

func (i *ImageObject) Data() interface{} {
	return i.data
}

func (i *ImageObject) Width() int {
	return i.data.Rect.Max.X
}

func (i *ImageObject) Height() int {
	return i.data.Rect.Max.Y
}

func loadImage(r Resource) (Image, error) {
	reader := r.reader
	if r.reader == nil {
		file, err := os.Open(r.url)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return &ImageObject{newm}, nil
}

func loadJson(r Resource) (string, error) {
	file, err := ioutil.ReadFile(r.url)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

type Assets struct {
	queue  []string
	cache  map[string]Image
	loads  int
	errors int
}

func NewAssets() *Assets {
	return &Assets{make([]string, 0), make(map[string]Image), 0, 0}
}

func (a *Assets) Image(path string) {
	a.queue = append(a.queue, path)
}

func (a *Assets) Get(path string) Image {
	return a.cache[path]
}

func (a *Assets) Load(onFinish func()) {
	if len(a.queue) == 0 {
		onFinish()
	} else {
		for _, path := range a.queue {
			img := LoadImage(path)
			a.cache[path] = img
		}
	}
}

func LoadImage(data interface{}) Image {
	var m image.Image

	switch data := data.(type) {
	default:
		log.Fatal("NewTexture needs a string or io.Reader")
	case string:
		file, err := os.Open(data)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case io.Reader:
		img, _, err := image.Decode(data)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case image.Image:
		m = data
	}

	b := m.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), m, b.Min, draw.Src)

	return &ImageObject{newm}
}

func Minimize() {
	fatalErr(window.Iconify())
}

func AppDir() string {
	dir, err := osext.ExecutableFolder()
	fatalErr(err)
	return dir
}

type JoystickState struct {
	Buttons []byte
	Axes    []float32
}

func Joysticks() map[int]JoystickState {
	states := make(map[int]JoystickState)
	for js := glfw.Joystick1; js <= glfw.JoystickLast; js++ {
		if !glfw.JoystickPresent(js) {
			continue
		}
		states[int(js)] = JoystickState{
			Buttons: glfw.GetJoystickButtons(js),
			Axes:    glfw.GetJoystickAxes(js),
		}
	}
	return states
}

// This is going away, don't use it.
func DumpJoysticks() {
	var present []glfw.Joystick
	for i := glfw.Joystick1; i <= glfw.JoystickLast; i++ {
		active := glfw.JoystickPresent(i)
		if active {
			present = append(present, i)
		}
	}

	go func() {
		for {
			for _, js := range present {
				log.Println(js, glfw.GetJoystickName(js), glfw.GetJoystickButtons(js))
				log.Println("   ", glfw.GetJoystickAxes(js))
			}
			time.Sleep(1 * time.Second)
		}
	}()
}
