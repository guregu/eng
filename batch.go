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
	Width() float32
	Height() float32
	View() (float32, float32, float32, float32)
}

type Batch struct {
	screen *ebiten.Image
}

func NewBatch(width, height float32) *Batch {
	batch := &Batch{
		screen: screen,
	}
	return batch
	// batch := new(Batch)
	// img, _ ebiten.NewImage(width, height)

	// batch.shader = LoadShader(batchVert, batchFrag)
	// batch.inPosition = gl.GetAttribLocation(batch.shader, "in_Position")
	// batch.inColor = gl.GetAttribLocation(batch.shader, "in_Color")
	// batch.inTexCoords = gl.GetAttribLocation(batch.shader, "in_TexCoords")
	// batch.ufProjection = gl.GetUniformLocation(batch.shader, "uf_Projection")

	// batch.vertices = make([]float32, 20*size)
	// batch.indices = make([]uint16, 6*size)

	// for i, j := 0, 0; i < size*6; i, j = i+6, j+4 {
	// 	batch.indices[i+0] = uint16(j + 0)
	// 	batch.indices[i+1] = uint16(j + 1)
	// 	batch.indices[i+2] = uint16(j + 2)
	// 	batch.indices[i+3] = uint16(j + 0)
	// 	batch.indices[i+4] = uint16(j + 2)
	// 	batch.indices[i+5] = uint16(j + 3)
	// }

	// batch.indexVBO = gl.CreateBuffer()
	// batch.vertexVBO = gl.CreateBuffer()

	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, batch.indexVBO)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, batch.indices, gl.STATIC_DRAW)

	// gl.BindBuffer(gl.ARRAY_BUFFER, batch.vertexVBO)
	// gl.BufferData(gl.ARRAY_BUFFER, batch.vertices, gl.DYNAMIC_DRAW)

	// gl.EnableVertexAttribArray(batch.inPosition)
	// gl.EnableVertexAttribArray(batch.inTexCoords)
	// gl.EnableVertexAttribArray(batch.inColor)

	// gl.VertexAttribPointer(batch.inPosition, 2, gl.FLOAT, false, 20, 0)
	// gl.VertexAttribPointer(batch.inTexCoords, 2, gl.FLOAT, false, 20, 8)
	// gl.VertexAttribPointer(batch.inColor, 4, gl.UNSIGNED_BYTE, true, 20, 16)

	// batch.projX = width / 2
	// batch.projY = height / 2

	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// return batch
}

func (b *Batch) Begin() {
	// if b.drawing {
	// 	log.Fatal("Batch.End() must be called first")
	// }
	// b.drawing = true
	// gl.UseProgram(b.shader)
}

func (b *Batch) End() {
	// b.screen.
	// 	if !b.drawing {
	// 		log.Fatal("Batch.Begin() must be called first")
	// 	}
	// 	if b.index > 0 {
	// 		b.flush()
	// 	}
	// 	b.drawing = false

	// 	b.lastTexture = nil
}

func (b *Batch) SetProjection(width, height float32) {
	fmt.Println("setproj", width, height)
	// TODO
	// b.projX = width / 2
	// b.projY = height / 2
}

func (b *Batch) Draw(r Drawable, x, y, originX, originY, scaleX, scaleY, rotation float32, colorpack uint32, transparency float32) {
	scr := screen
	if scr == nil {
		log.Println("scr is nil")
		return
	}
	opt := &ebiten.DrawImageOptions{}

	opt.GeoM.Scale(float64(scaleX), float64(scaleY))
	opt.GeoM.Translate(float64(x), float64(y))
	opt.GeoM.Rotate(float64(rotation))

	blue := float64(uint8(colorpack & 0x000000FF))         // 10
	green := float64(uint8((colorpack & 0x0000FF00) >> 8)) // 154
	red := float64(uint8((colorpack & 0x00FF0000) >> 16))  // 0
	alpha := float64(transparency)

	opt.ColorM.Scale(red/255.0, green/255.0, blue/255.0, alpha)

	tex := r.Texture()
	scr.DrawImage(tex, opt)
}
