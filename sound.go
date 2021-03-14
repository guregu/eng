package engi

import (
	// "log"
	// "os"
	"bytes"
	"context"
	// "io"
	// "encoding/binary"
	"io"
	"math"
	"time"
	// "unsafe"
	// "azul3d.org/audio.v1"
	// "github.com/guregu/native-al"
	// "golang.org/x/net/context"
	// _ "azul3d.org/audio/wav.v1"
	// _ "github.com/guregu/audio-flac"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mewkiz/flac"

	"log"
)

var (
	// audioDevice  *al.Device
	audioSources []uint32
	audioContext context.Context
	audioCancel  func()
	muted        bool
	gain         float64
)

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
	source  uint32
	buffers []uint32
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
	if s == nil {
		return
	}
	s.player.SetVolume(gain)
	s.player.Play()
	// s.looping = true
	// s.play()
}

func (s *Music) play() {
	// s.context, s.cancel = context.WithCancel(audioContext)
	// s.source = nextAvailableSource()
	// s.reset()
	// audioDevice.Sourcei(s.source, al.LOOPING, al.FALSE)
	// // fill all buffers first
	// for _, buf := range s.buffers {
	// 	s.fill(buf)
	// }

	// audioDevice.SourceQueueBuffers(s.source, s.buffers)
	// go s.run()

	// audioDevice.SourcePlay(s.source)
}

func (s *Music) reset() {
	// flac pkg can't seek yet so ghetto seek to 0
	// s.file.Seek(0, os.SEEK_SET)
	// var err error
	// s.decoder, _, err = audio.NewDecoder(s.file)
	// if err != nil {
	// 	panic(err)
	// }
}

func (s *Music) unqueue() []uint32 {
	// var processed int32

	// audioDevice.GetSourcei(s.source, al.BUFFERS_PROCESSED, &processed)
	// if processed == 0 {
	// 	return nil
	// }
	// available := make([]uint32, processed)
	// audioDevice.SourceUnqueueBuffers(s.source, available)
	// return available
	return nil
}

// func (s *Music) fill(buffer uint32) error {
// 	config := s.decoder.Config()
// 	bufSize := config.SampleRate
// 	samples := make(audio.PCM16Samples, bufSize)

// 	read, err := s.decoder.Read(samples)
// 	if err != nil && err != audio.EOS {
// 		return err
// 	}

// 	if read > 0 {
// 		data := []audio.PCM16(samples[:read])
// 		if config.Channels == 1 {
// 			audioDevice.BufferData(buffer, al.FORMAT_MONO16, unsafe.Pointer(&data[0]), int32(int(unsafe.Sizeof(data[0]))*read), int32(config.SampleRate))
// 		} else {
// 			audioDevice.BufferData(buffer, al.FORMAT_STEREO16, unsafe.Pointer(&data[0]), int32(int(unsafe.Sizeof(data[0]))*read), int32(config.SampleRate))
// 		}
// 	}

// 	return err
// }

func (s *Music) Delete() {
	if s == nil {
		return
	}
	s.player.Close()
	// if s.cancel != nil {
	// 	s.cancel()
	// }
	// s.file.Close()
	// audioDevice.DeleteBuffers(buffersPerStream, &s.buffers[0])
}

func (s *Music) Playing() bool {
	if s == nil || s.player == nil {
		return false
	}
	return s.player.IsPlaying()
}

func (s *Music) Stop() {
	if s == nil {
		return
	}
	if s.cancel != nil {
		s.cancel()
	}
	// audioDevice.SourceStop(s.source)
	// s.unqueue()
	s.source = 0
}

type flacStream struct {
	stream *flac.Stream
}

func (fs *flacStream) Read(b []byte) (n int, err error) {
	// defer func() {
	b = b[:n]
	// }()
	// b = b[:0]

more:
	// frame, err := fs.stream.ParseNext()
	// if err != nil {
	// 	if err == io.EOF {
	// 		return n, io.EOF // TODO
	// 	}
	// 	return n, err
	// }
	// out := bytes.NewBuffer(b)
	_ = bytes.Compare
	// var frames int
	// println(fs.stream.)
	// log.Printf("%+v\n", fs.stream.Info)
	// for i := 0; i < frame.Subframes[0].NSamples; i++ {
	// 	for x, subframe := range frame.Subframes {
	// 		if x > 1 {
	// 			// continue
	// 		}
	// 		samp := uint32(subframe.Samples[i])
	// 		const max = 32767
	// 		q := int16(samp)
	_ = math.Pi
	_ = log.Println
	// 		// q := int16(math.Sin(2*math.Pi*float64(samp)/2) * 0.3 * max)
	// 		b[n] = byte(q)
	// 		n++
	// 		b[n] = byte(q >> 8)
	// 		n++

	// 		// for j := uint32(2); j < 4; j++ {
	// 		// 	f := byte((samp >> (8 * j)) & 0xff)

	// 		// }
	// 		if n == len(b) {
	// 			return n, nil
	// 		}
	// 		// data = append(data, int(subframe.Samples[i]))
	// 	}
	// }

	// out:
	// 	for x := 0; n+4 < len(b); x++ {
	// 		for _, sub := range frame.Subframes {
	// 			if frames > len(sub.Samples) {
	// 				break out
	// 			}
	// 			// sub.
	// 			samp := sub.Samples[frames]
	// 			frames++

	// 			for i := uint32(0); i < 4; i++ {
	// 				b[n] = byte((samp >> (8 * i)) & 0xff)
	// 				n++
	// 			}
	// 		}
	// 	}
	if n+4 < len(b) {
		goto more
	}
	return n, nil
}

func loadMusic(r Resource) (*Music, error) {
	return nil, nil

	f, err := ebitenutil.OpenFile(r.url)
	if err != nil {
		return nil, err
	}
	stream, err := flac.New(f)
	if err != nil {
		return nil, err
	}
	fs := &flacStream{
		stream: stream,
	}
	p, err := audio.NewPlayer(acx, fs)
	if err != nil {
		return nil, err
	}
	return &Music{
		player: p,
	}, nil
	// audio.
	// wav.
	return nil, nil
	// func readSoundFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {

	// file, err := os.Open(r.url)
	// if err != nil {
	// 	return nil, err
	// }
	// // audio.
	// // audio.
	// s := Music{
	// 	// buffers: make([]uint32, buffersPerStream),
	// 	// file:    file,
	// }

	// // audioDevice.GenBuffers(buffersPerStream, &s.buffers[0])

	// return &s, nil
}

type SFX struct {
	buf        []byte
	player     *audio.Player
	source     uint32
	buffer     uint32
	duration   time.Duration // seconds
	sampleRate int
	looping    bool
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

var acx = audio.NewContext(44100)

func loadSFX(res Resource) (*SFX, error) {
	r, err := ebitenutil.OpenFile(res.url)
	if err != nil {
		return nil, err
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
