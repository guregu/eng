// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Action int

// type Key int
type Modifier int

type Key = ebiten.Key

var (
	MOVE    = Action(0)
	PRESS   = Action(1)
	RELEASE = Action(2)
	SHIFT   = Modifier(0x0001)
	CONTROL = Modifier(0x0002)
	ALT     = Modifier(0x0004)
	SUPER   = Modifier(0x0008)
)

var (
	Dash        = ebiten.KeyMinus
	Apostrophe  = ebiten.KeyApostrophe
	Semicolon   = ebiten.KeySemicolon
	Equals      = ebiten.KeyEqual
	Comma       = ebiten.KeyComma
	Period      = ebiten.KeyPeriod
	Slash       = ebiten.KeySlash
	Backslash   = ebiten.KeyBackslash
	Backspace   = ebiten.KeyBackspace
	Tab         = ebiten.KeyTab
	CapsLock    = ebiten.KeyCapsLock
	Space       = ebiten.KeySpace
	Enter       = ebiten.KeyEnter
	Escape      = ebiten.KeyEscape
	Insert      = ebiten.KeyInsert
	PrintScreen = ebiten.KeyPrintScreen
	Delete      = ebiten.KeyDelete
	PageUp      = ebiten.KeyPageUp
	PageDown    = ebiten.KeyPageDown
	Home        = ebiten.KeyHome
	End         = ebiten.KeyEnd
	Pause       = ebiten.KeyPause
	ScrollLock  = ebiten.KeyScrollLock
	ArrowLeft   = ebiten.KeyLeft
	ArrowRight  = ebiten.KeyRight
	ArrowDown   = ebiten.KeyDown
	ArrowUp     = ebiten.KeyUp
	LeftBracket = ebiten.KeyLeftBracket
	LeftShift   = ebiten.KeyShift
	LeftControl = ebiten.KeyControl
	// LeftSuper    = ebiten.KeySuper
	LeftAlt      = ebiten.KeyAlt
	RightBracket = ebiten.KeyRightBracket
	RightShift   = ebiten.KeyShift
	RightControl = ebiten.KeyControl
	// RightSuper   = ebiten.Super
	RightAlt = ebiten.KeyAlt
	Zero     = ebiten.Key0
	One      = ebiten.Key1
	Two      = ebiten.Key2
	Three    = ebiten.Key3
	Four     = ebiten.Key4
	Five     = ebiten.Key5
	Six      = ebiten.Key6
	Seven    = ebiten.Key7
	Eight    = ebiten.Key8
	Nine     = ebiten.Key9
	F1       = ebiten.KeyF1
	F2       = ebiten.KeyF2
	F3       = ebiten.KeyF3
	F4       = ebiten.KeyF4
	F5       = ebiten.KeyF5
	F6       = ebiten.KeyF6
	F7       = ebiten.KeyF7
	F8       = ebiten.KeyF8
	F9       = ebiten.KeyF9
	F10      = ebiten.KeyF10
	F11      = ebiten.KeyF11
	F12      = ebiten.KeyF12
	A        = ebiten.KeyA
	B        = ebiten.KeyB
	C        = ebiten.KeyC
	D        = ebiten.KeyD
	E        = ebiten.KeyE
	F        = ebiten.KeyF
	G        = ebiten.KeyG
	H        = ebiten.KeyH
	I        = ebiten.KeyI
	J        = ebiten.KeyJ
	K        = ebiten.KeyK
	L        = ebiten.KeyL
	M        = ebiten.KeyM
	N        = ebiten.KeyN
	O        = ebiten.KeyO
	P        = ebiten.KeyP
	Q        = ebiten.KeyQ
	R        = ebiten.KeyR
	S        = ebiten.KeyS
	T        = ebiten.KeyT
	U        = ebiten.KeyU
	V        = ebiten.KeyV
	W        = ebiten.KeyW
	X        = ebiten.KeyX
	Y        = ebiten.KeyY
	Z        = ebiten.KeyZ
	NumLock  = ebiten.KeyNumLock
	// NumMultiply = ebiten.KeyNumMultiply
	// NumDivide   = ebiten.KeyNumDivide
	// NumAdd      = ebiten.KeyNumAdd
	// NumSubtract = ebiten.KeyNumSubtract
	// NumZero     = ebiten.KeyNumZero
	// NumOne      = ebiten.KeyNumOne
	// NumTwo      = ebiten.KeyNumTwo
	// NumThree    = ebiten.KeyNumThree
	// NumFour     = ebiten.KeyNumFour
	// NumFive     = ebiten.KeyNumFive
	// NumSix      = ebiten.KeyNumSix
	// NumSeven    = ebiten.KeyNumSeven
	// NumEight    = ebiten.KeyNumEight
	// NumNine     = ebiten.KeyNumNine
	// NumDecimal  = ebiten.KeyNumDecimal
	// NumEnter    = ebiten.KeyNumEnter
)
