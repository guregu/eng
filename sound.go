package engi

import (
	"log"
	// "os"
	"bytes"
	"context"
	// "io"
	// "encoding/binary"
	"io"
	// "log"
	// "math"
	"time"
	// "unsafe"
	// "azul3d.org/audio.v1"
	// "github.com/guregu/native-al"
	// "golang.org/x/net/context"
	// _ "azul3d.org/audio/wav.v1"
	// _ "github.com/guregu/audio-flac"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	// audioDevice  *al.Device
	audioSources []uint32
	// audioContext context.Context
	audioCancel func()
	muted       bool
	gain        float64
)
var acx = audio.NewContext(44100)

const (
	audioSourcesCount = 31
	buffersPerStream  = 3
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
	_ = g

	// for _, src := range audioSources {
	// 	audioDevice.Sourcef(src, al.GAIN, g)
	// }
}

func Muted() bool {
	return muted
}

func SetGain(g float64) {
	if g != gain {
		gain = g

		if muted {
			return
		}

		// for _, src := range audioSources {
		// 	audioDevice.Sourcef(src, al.GAIN, g)
		// }
	}
}

type Music struct {
	player  *audio.Player
	src     io.ReadSeekCloser
	looping bool
	context context.Context
	cancel  func()

	// decoder audio.Decoder
	// file    *os.File
}

func (s *Music) Play() {
	if s == nil {
		return
	}
	s.player.SetVolume(gain)
	s.player.Play()
	// s.looping = false
	// s.play()
}

func (s *Music) PlayAt(mult float64) {
	s.Play()
}

func (s *Music) Loop() {
	s.Play()
	// if s == nil {
	// 	return
	// }
	// s.looping = true
	// s.player.SetVolume(gain)
	// s.player.Play()
}

func (s *Music) Delete() {
	if s == nil {
		return
	}
	s.player.Close()
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

	log.Println("loaded", r)

	return &Music{
		player: player,
	}, nil

	// f, err := ebitenutil.OpenFile(r.url)
	// if err != nil {
	// 	return nil, err
	// }
	// stream, err := flac.New(f)
	// if err != nil {
	// 	return nil, err
	// }
	// fs := &flacStream{
	// 	stream: stream,
	// }
	// p, err := audio.NewPlayer(acx, fs)
	// if err != nil {
	// 	return nil, err
	// }
	// return &Music{
	// 	player: p,
	// }, nil
	// // audio.
	// // wav.
	// return nil, nil
	// // func readSoundFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {

	// // file, err := os.Open(r.url)
	// // if err != nil {
	// // 	return nil, err
	// // }
	// // // audio.
	// // // audio.
	// // s := Music{
	// // 	// buffers: make([]uint32, buffersPerStream),
	// // 	// file:    file,
	// // }

	// // // audioDevice.GenBuffers(buffersPerStream, &s.buffers[0])

	// return &s, nil
}

type SFX struct {
	buf      []byte
	player   *audio.Player
	duration time.Duration // seconds
	looping  bool
}

// func (s *SFX) bind() {
// 	s.source = nextAvailableSource()
// 	audioDevice.Sourcei(s.source, al.BUFFER, int32(s.buffer))
// 	if s.looping {
// 		audioDevice.Sourcei(s.source, al.LOOPING, al.TRUE)
// 	} else {
// 		audioDevice.Sourcei(s.source, al.LOOPING, al.FALSE)
// 	}
// }

func (s *SFX) Play() {
	s.play(gain)
}

func (s *SFX) PlayAt(gainMultiplier float64) {
	gain := gain
	gain *= gainMultiplier
	s.play(gain)
}

func (s *SFX) play(gain float64) {
	if s == nil {
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

// func readSFXFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, audio.Config{}, 0, err
// 	}
// 	defer file.Close()
// 	fi, err := file.Stat()
// 	if err != nil {
// 		return nil, audio.Config{}, 0, err
// 	}

// 	decoder, _, err := audio.NewDecoder(file)
// 	if err != nil {
// 		return nil, audio.Config{}, 0, err
// 	}

// 	config = decoder.Config()

// 	// Convert everything to 16-bit samples
// 	bufSize := int(fi.Size())
// 	samples = make(audio.PCM16Samples, 0, bufSize)

// 	// TODO: surely there is a better way to do this
// 	var read int
// 	buf := make(audio.PCM16Samples, 1024*1024)
// 	err = nil
// 	for err != audio.EOS {
// 		var r int
// 		r, err = decoder.Read(buf)
// 		if err != nil && err != audio.EOS {
// 			return nil, audio.Config{}, 0, err
// 		}
// 		read += r
// 		samples = append(samples, buf[:r]...)
// 	}

// 	secs := 1 / float64(config.SampleRate) * float64(read)
// 	duration = time.Duration(float64(time.Second) * secs)
// 	return []audio.PCM16(samples)[:read], config, duration, nil
// }

// func cleanupAudio() {
// 	audioCancel()
// 	audioDevice.DeleteSources(int32(len(audioSources)), &audioSources[0])
// 	for _, s := range Files.sfx {
// 		s.Delete()
// 	}
// 	for _, s := range Files.music {
// 		s.Delete()
// 	}
// 	if audioDevice != nil {
// 		audioDevice.Close()
// 	}
// }

// func init() {
// 	al.SetErrorHandler(func(err error) {
// 		log.Println("[audio]", err)
// 	})
// }

var _ Sound = (*Music)(nil)
var _ Sound = (*SFX)(nil)
