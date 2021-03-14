// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	// "image/color"
	"log"
	// "math"
	"fmt"

	// webgl "engo.io/gl"

	_ "github.com/davecgh/go-spew/spew"

	"github.com/hajimehoshi/ebiten/v2"
)

const size = 10000

type Drawable interface {
	Texture() *ebiten.Image
	Width() float64
	Height() float64
	View() (float64, float64, float64, float64)
}

type Batch struct {
	screen *ebiten.Image
}

func NewBatch(width, height float64) *Batch {
	batch := &Batch{
		screen: screen,
	}
	return batch
}

func (b *Batch) Begin() {}

func (b *Batch) End() {}

func (b *Batch) SetProjection(width, height float64) {
	fmt.Println("setproj", width, height)
	// TODO
	// b.projX = width / 2
	// b.projY = height / 2
}

func (b *Batch) Draw(r Drawable, x, y, originX, originY, scaleX, scaleY, rotation float64, colorpack uint32, transparency float64) {
	scr := screen
	if scr == nil {
		log.Println("scr is nil")
		return
	}
	opt := &ebiten.DrawImageOptions{}

	opt.GeoM.Scale(scaleX, scaleY)
	opt.GeoM.Translate(x, y)
	opt.GeoM.Rotate(rotation)

	if colorpack != 0xffffff || transparency != 1 {
		blue := float64(uint8(colorpack & 0x000000FF))         // 10
		green := float64(uint8((colorpack & 0x0000FF00) >> 8)) // 154
		red := float64(uint8((colorpack & 0x00FF0000) >> 16))  // 0
		alpha := float64(transparency)
		opt.ColorM.Scale(red/255.0, green/255.0, blue/255.0, alpha)
	}

	tex := r.Texture()
	scr.DrawImage(tex, opt)
}

func (b *Batch) Draw2(img *ebiten.Image, opt *ebiten.DrawImageOptions) {
	scr := screen
	if scr == nil {
		log.Println("scr is nil")
		return
	}
	scr.DrawImage(img, opt)
}
