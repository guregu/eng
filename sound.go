package engi

import (
	"bytes"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	acx     = audio.NewContext(44100)
	playing = make(map[*Music]struct{})

	muted bool
	gain  float64
)

type Sound interface {
	Play()
	PlayAt(mult float64)
	Loop()
	Stop()
	Playing() bool
}

func ToggleMute() {
	muted = !muted
	g := gain
	if muted {
		g = 0
	}
	applyGain(g)
}

func Muted() bool {
	return muted
}

func SetGain(g float64) {
	gain = g
	applyGain(g)
}

func applyGain(g float64) {
	for m := range playing {
		if m.player != nil {
			m.player.SetVolume(g)
		}
	}
}

func volume() float64 {
	if muted {
		return 0
	}
	return gain
}

type Music struct {
	player  *audio.Player
	looping bool
}

func (s *Music) Play() {
	if s == nil {
		return
	}
	s.player.SetVolume(volume())
	s.player.Play()
	playing[s] = struct{}{}
}

func (s *Music) PlayAt(mult float64) {
	s.Play()
}

func (s *Music) Loop() {
	// we just assume all Music is looping...
	s.Play()
}

func (s *Music) Delete() {
	if s == nil {
		return
	}
	s.player.Close()
	delete(playing, s)
}

func (s *Music) Playing() bool {
	if s == nil {
		return false
	}
	return s.player.IsPlaying()
}

func (s *Music) Stop() {
	if s == nil {
		return
	}
	s.player.Pause()
	s.player.Rewind()
	delete(playing, s)
}

func loadMusic(r Resource) (*Music, error) {
	var f ebitenutil.ReadSeekCloser
	var err error
	if rsc, ok := r.reader.(ebitenutil.ReadSeekCloser); ok {
		f = rsc
	} else if r.url != "" {
		f, err = ebitenutil.OpenFile(r.url)
	} else {
		f, err = ebitenutil.OpenFile("data/" + r.name + ".mp3")
	}

	stream, err := mp3.Decode(acx, f)
	if err != nil {
		return nil, err
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())

	player, err := audio.NewPlayer(acx, loop)
	if err != nil {
		return nil, err
	}

	return &Music{
		player: player,
	}, nil
}

type SFX struct {
	buf      []byte
	player   *audio.Player
	duration time.Duration // seconds
	looping  bool
}

func (s *SFX) Play() {
	s.play(volume())
}

func (s *SFX) PlayAt(gainMultiplier float64) {
	gain := volume()
	gain *= gainMultiplier
	s.play(gain)
}

func (s *SFX) play(gain float64) {
	if s == nil {
		return
	}
	if gain < 0.05 {
		return
	}
	player := audio.NewPlayerFromBytes(acx, s.buf)
	player.SetVolume(gain)
	player.Play()
	s.player = player
}

func (s *SFX) Loop() {
	// s.looping = true
	// s.bind()
	// audioDevice.SourcePlay(s.source)
}

func (s *SFX) Stop() {
	if s == nil {
		return
	}
	if s.player != nil {
		s.player.Close()
		s.player = nil
	}
}

func (s *SFX) Delete() {
	// audioDevice.DeleteBuffers(1, &s.buffer)
}

func (s *SFX) Playing() bool {
	if s == nil {
		return false
	}
	if s.player != nil {
		return s.Playing()
	}
	return false
}

func (s *SFX) Duration() time.Duration {
	if s == nil {
		return 0
	}
	return s.duration
}

func loadSFX(res Resource) (*SFX, error) {
	r, _ := res.reader.(io.ReadSeeker)
	if r == nil {
		var err error
		r, err = ebitenutil.OpenFile(res.url)
		if err != nil {
			return nil, err
		}
	}
	x, err := wav.Decode(acx, r)
	if err != nil {
		return nil, err
	}
	x.Length()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, x); err != nil {
		return nil, err
	}
	return &SFX{
		buf: buf.Bytes(),
	}, nil
}

var _ Sound = (*Music)(nil)
var _ Sound = (*SFX)(nil)
