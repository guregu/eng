// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"image"
	"io"
	// "math"
	"path"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/davecgh/go-spew/spew"
	// webgl "engo.io/gl"
)

type Resource struct {
	kind   string
	name   string
	url    string
	reader io.Reader
}

type Loader struct {
	resources []Resource
	loaded    map[Resource]struct{}
	images    map[string]*Texture
	jsons     map[string]string
	sfx       map[string]*SFX
	music     map[string]*Music
}

func NewLoader() *Loader {
	return &Loader{
		resources: make([]Resource, 1),
		loaded:    make(map[Resource]struct{}),
		images:    make(map[string]*Texture),
		jsons:     make(map[string]string),
		sfx:       make(map[string]*SFX),
		music:     make(map[string]*Music),
	}
}

func (l *Loader) Add(name, url string) {
	kind := path.Ext(url)[1:]
	url = filepath.FromSlash(url)
	l.resources = append(l.resources, Resource{
		kind: kind,
		name: name,
		url:  url,
	})
}

func (l *Loader) AddReader(name, kind string, r io.Reader) {
	l.resources = append(l.resources, Resource{
		kind:   kind,
		name:   name,
		reader: r,
	})
}

func (l *Loader) Image(name string) *Texture {
	return l.images[name]
}

func (l *Loader) Json(name string) string {
	return l.jsons[name]
}

func (l *Loader) Sound(name string) Sound {
	if sfx, ok := l.sfx[name]; ok {
		return sfx
	}
	return l.music[name]
}

func (l *Loader) SFX(name string) *SFX {
	return l.sfx[name]
}

func (l *Loader) Music(name string) *Music {
	return l.music[name]
}

func (l *Loader) Load(onFinish func()) {
	for _, r := range l.resources {
		if _, loaded := l.loaded[r]; loaded {
			// don't load stuff twice
			continue
		}
		switch r.kind {
		case "png":
			data, err := loadImage(r)
			fatalErr(err)
			l.images[r.name] = NewTexture(data)
			l.loaded[r] = struct{}{}
			// spew.Dump(data, err)
		case "json":
			data, err := loadJson(r)
			fatalErr(err)
			l.jsons[r.name] = data
			l.loaded[r] = struct{}{}
		case "wav" /*, "flac-sfx"*/ :
			data, err := loadSFX(r)
			fatalErr(err)
			l.sfx[r.name] = data
			l.loaded[r] = struct{}{}
		case "flac":
			data, err := loadMusic(r)
			fatalErr(err)
			l.music[r.name] = data
			l.loaded[r] = struct{}{}
		}
	}
	onFinish()
}

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

type Region struct {
	texture *Texture
	img     *ebiten.Image
	x, y    int
	w, h    int
}

func NewRegion(texture *Texture, x, y, w, h int) *Region {
	// TODO
	sub := texture.img.SubImage(image.Rect(x, y, x+w, y+h))
	// img := ebiten.NewImageFromImage(sub)
	return &Region{
		// texture: texture,
		// img:     img,
		img: sub.(*ebiten.Image),
		x:   x,
		y:   y,
		w:   w,
		h:   h,
	}
}

func (r *Region) Width() float64 {
	return float64(r.w)
}

func (r *Region) Height() float64 {
	return float64(r.h)
}

func (r *Region) Texture() *ebiten.Image {
	return r.img
}

func (r *Region) View() (float64, float64, float64, float64) {
	// return r.u, r.v, r.u2, r.v2
	return 0.0, 0.0, 1.0, 1.0
}

type Texture struct {
	img    *ebiten.Image
	width  int
	height int
}

func NewTexture(img Image) *Texture {
	tex := ebiten.NewImageFromImage(img.Data().(*image.NRGBA))
	return &Texture{tex, img.Width(), img.Height()}
}

// Width returns the width of the texture.
func (t *Texture) Width() float64 {
	return float64(t.width)
}

// Height returns the height of the texture.
func (t *Texture) Height() float64 {
	return float64(t.height)
}

func (t *Texture) Texture() *ebiten.Image {
	return t.img
}

func (r *Texture) View() (float64, float64, float64, float64) {
	return 0.0, 0.0, 1.0, 1.0
}

type Point struct {
	X, Y float64
}

func (p *Point) Set(x, y float64) {
	p.X = x
	p.Y = y
}

func (p *Point) SetTo(v float64) {
	p.X = v
	p.Y = v
}

type Sprite struct {
	Position *Point
	Scale    *Point
	Anchor   *Point
	Rotation float64
	Color    uint32
	Alpha    float64
	Region   *Region
}

func NewSprite(region *Region, x, y float64) *Sprite {
	return &Sprite{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Anchor:   &Point{0, 0},
		Rotation: 0,
		Color:    0xffffff,
		Alpha:    1,
		Region:   region,
	}
}

func (s *Sprite) Render(batch *Batch) {
	batch.Draw(s.Region, s.Position.X, s.Position.Y, s.Anchor.X, s.Anchor.Y, s.Scale.X, s.Scale.Y, s.Rotation, s.Color, s.Alpha)
}
