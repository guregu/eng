// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"math"
	"time"

	// azclock "github.com/azul3d/engine/gfx/clock"
)

type Clock struct {
	// *azclock.Clock
	elapsed float64
	delta   float64
	fps     float64
	frames  uint64
	start   time.Time
	frame   time.Time
}

func NewClock() *Clock {
	clock := new(Clock)
	clock.start = time.Now()
	clock.Tick()
	return clock

	// clock := new(Clock)
	// clock.Clock = azclock.New()
	// // clock.SetMaxFrameRate(75)
	// clock.Tick()
	// return clock
}

func (c *Clock) Tick() {
	// c.Clock.Tick()

	now := time.Now()
	c.frames += 1
	c.delta = now.Sub(c.frame).Seconds()
	c.elapsed += c.delta
	c.frame = now

	if c.elapsed >= 1 {
		c.fps = float64(c.frames)
		c.elapsed = math.Mod(c.elapsed, 1)
		c.frames = 0
	}
}

func (c *Clock) Delta() float32 {
	// return float32(c.Clock.Dt())
	return float32(c.delta)
}

func (c *Clock) Fps() float32 {
	// return float32(c.Clock.FrameRate())
	return float32(c.fps)
}

func (c *Clock) Time() float32 {
	// return float32(c.Clock.Time().Seconds())
	return float32(time.Now().Sub(c.start).Seconds())
}
