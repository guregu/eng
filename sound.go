package engi

import (
	"log"
	"os"
	"time"
	"unsafe"

	"azul3d.org/audio.v1"
	"azul3d.org/native/al.v1"
	"golang.org/x/net/context"

	_ "azul3d.org/audio/flac.dev"
	_ "azul3d.org/audio/wav.v1"
)

var (
	audioDevice  *al.Device
	audioSources []uint32
	audioContext context.Context
	audioCancel  func()
	muted        bool
	gain         float32
)

const (
	audioSourcesCount = 31
	buffersPerStream  = 3
)

type Sound interface {
	Play()
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

	for _, src := range audioSources {
		audioDevice.Sourcef(src, al.GAIN, g)
	}
}

func Muted() bool {
	return muted
}

func SetGain(g float32) {
	if g != gain {
		gain = g

		if muted {
			return
		}

		for _, src := range audioSources {
			audioDevice.Sourcef(src, al.GAIN, g)
		}
	}
}

type Music struct {
	source  uint32
	buffers []uint32
	looping bool
	context context.Context
	cancel  func()

	decoder audio.Decoder
	file    *os.File
}

func (s *Music) Play() {
	s.looping = false
	s.play()
}

func (s *Music) Loop() {
	s.looping = true
	s.play()
}

func (s *Music) play() {
	s.context, s.cancel = context.WithCancel(audioContext)
	s.source = nextAvailableSource()
	s.reset()
	audioDevice.Sourcei(s.source, al.LOOPING, al.FALSE)
	// fill all buffers first
	for _, buf := range s.buffers {
		s.fill(buf)
	}

	audioDevice.SourceQueueBuffers(s.source, s.buffers)
	go s.run()

	audioDevice.SourcePlay(s.source)
}

func (s *Music) reset() {
	// flac pkg can't seek yet so ghetto seek to 0
	s.file.Seek(0, os.SEEK_SET)
	var err error
	s.decoder, _, err = audio.NewDecoder(s.file)
	if err != nil {
		panic(err)
	}
}

func (s *Music) run() {
	for {
		select {
		case <-s.context.Done():
			return
		default:
			processed := s.unqueue()
			if len(processed) == 0 {
				continue
			}

			for _, buf := range processed {
				err := s.fill(buf)
				switch {
				case err == audio.EOS:
					if s.looping {
						// time.Sleep(500 * time.Millisecond)
						// start over
						s.reset()
					}
				case err != nil:
					panic(err)
				}
				audioDevice.SourceQueueBuffers(s.source, []uint32{buf})
				if err == audio.EOS && !s.looping {
					return
				}
			}
		}
	}
}

func (s *Music) unqueue() []uint32 {
	var processed int32

	audioDevice.GetSourcei(s.source, al.BUFFERS_PROCESSED, &processed)
	if processed == 0 {
		return nil
	}
	available := make([]uint32, processed)
	audioDevice.SourceUnqueueBuffers(s.source, available)
	return available
}

func (s *Music) fill(buffer uint32) error {
	config := s.decoder.Config()
	bufSize := config.SampleRate
	samples := make(audio.PCM16Samples, bufSize)

	read, err := s.decoder.Read(samples)
	if err != nil && err != audio.EOS {
		return err
	}

	if read > 0 {
		data := []audio.PCM16(samples[:read])
		if config.Channels == 1 {
			audioDevice.BufferData(buffer, al.FORMAT_MONO16, unsafe.Pointer(&data[0]), int32(int(unsafe.Sizeof(data[0]))*read), int32(config.SampleRate))
		} else {
			audioDevice.BufferData(buffer, al.FORMAT_STEREO16, unsafe.Pointer(&data[0]), int32(int(unsafe.Sizeof(data[0]))*read), int32(config.SampleRate))
		}
	}

	return err
}

func (s *Music) Delete() {
	if s.cancel != nil {
		s.cancel()
	}
	s.file.Close()
	audioDevice.DeleteBuffers(buffersPerStream, &s.buffers[0])
}

func (s *Music) Playing() bool {
	if s.source == 0 {
		return false
	}
	var state int32
	audioDevice.GetSourcei(s.source, al.SOURCE_STATE, &state)
	return state == al.PLAYING
}

func (s *Music) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	audioDevice.SourceStop(s.source)
	s.unqueue()
	s.source = 0
}

func loadMusic(r Resource) (*Music, error) {
	// func readSoundFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {
	file, err := os.Open(r.url)
	if err != nil {
		return nil, err
	}

	s := Music{
		buffers: make([]uint32, buffersPerStream),
		file:    file,
	}

	audioDevice.GenBuffers(buffersPerStream, &s.buffers[0])

	return &s, nil
}

type SFX struct {
	source     uint32
	buffer     uint32
	duration   time.Duration // seconds
	sampleRate int
	looping    bool
}

func (s *SFX) bind() {
	s.source = nextAvailableSource()
	audioDevice.Sourcei(s.source, al.BUFFER, int32(s.buffer))
	if s.looping {
		audioDevice.Sourcei(s.source, al.LOOPING, al.TRUE)
	} else {
		audioDevice.Sourcei(s.source, al.LOOPING, al.FALSE)
	}
}

func (s *SFX) Play() {
	s.looping = false
	s.bind()
	audioDevice.SourcePlay(s.source)
}

func (s *SFX) Loop() {
	s.looping = true
	s.bind()
	audioDevice.SourcePlay(s.source)
}

func (s *SFX) Stop() {
	audioDevice.SourceStop(s.source)
	s.unqueue()
}

func (s *SFX) unqueue() []uint32 {
	var processed int32
	audioDevice.GetSourcei(s.source, al.BUFFERS_PROCESSED, &processed)
	if processed == 0 {
		return nil
	}
	available := make([]uint32, processed)
	audioDevice.SourceUnqueueBuffers(s.source, available)
	return available
}

func (s *SFX) Delete() {
	audioDevice.DeleteBuffers(1, &s.buffer)
}

func (s *SFX) Playing() bool {
	if s.source == 0 {
		return false
	}

	var state int32
	audioDevice.GetSourcei(s.source, al.SOURCE_STATE, &state)
	return state == al.PLAYING
}

func (s *SFX) Duration() time.Duration {
	return s.duration
}

func setupAudio() {
	var err error
	audioDevice, err = al.OpenDevice("", nil)
	fatalErr(err)

	audioSources = make([]uint32, audioSourcesCount)
	audioDevice.GenSources(audioSourcesCount, &audioSources[0])

	audioContext, audioCancel = context.WithCancel(context.Background())
}

func nextAvailableSource() uint32 {
	// find unused source
	for _, source := range audioSources {
		var state int32
		audioDevice.GetSourcei(source, al.SOURCE_STATE, &state)
		if state != al.PLAYING {
			audioDevice.Sourcei(source, al.BUFFER, 0)
			return source
		}
	}

	// no free sounds. find non-looping one and cut it short
	for _, source := range audioSources {
		var looping int32
		audioDevice.GetSourcei(source, al.LOOPING, &looping)
		if looping != al.TRUE {
			audioDevice.SourceStop(source)
			return source
		}
	}

	// give up, take the last one
	source := audioSources[len(audioSources)-1]
	audioDevice.SourceStop(source)
	return source
}

func loadSFX(r Resource) (*SFX, error) {
	samples, config, duration, err := readSFXFile(r.url)
	if err != nil {
		return nil, err
	}

	s := SFX{
		duration:   duration,
		sampleRate: config.SampleRate,
	}
	audioDevice.GenBuffers(1, &s.buffer)
	if config.Channels == 1 {
		audioDevice.BufferData(s.buffer, al.FORMAT_MONO16, unsafe.Pointer(&samples[0]), int32(int(unsafe.Sizeof(samples[0]))*len(samples)), int32(config.SampleRate))
	} else {
		audioDevice.BufferData(s.buffer, al.FORMAT_STEREO16, unsafe.Pointer(&samples[0]), int32(int(unsafe.Sizeof(samples[0]))*len(samples)), int32(config.SampleRate))
	}
	return &s, err
}

func readSFXFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, audio.Config{}, 0, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, audio.Config{}, 0, err
	}

	decoder, _, err := audio.NewDecoder(file)
	if err != nil {
		return nil, audio.Config{}, 0, err
	}

	config = decoder.Config()

	// Convert everything to 16-bit samples
	bufSize := int(fi.Size())
	samples = make(audio.PCM16Samples, 0, bufSize)

	// TODO: surely there is a better way to do this
	var read int
	buf := make(audio.PCM16Samples, 1024*1024)
	err = nil
	for err != audio.EOS {
		var r int
		r, err = decoder.Read(buf)
		if err != nil && err != audio.EOS {
			return nil, audio.Config{}, 0, err
		}
		read += r
		samples = append(samples, buf[:r]...)
	}

	secs := 1 / float64(config.SampleRate) * float64(read)
	duration = time.Duration(float64(time.Second) * secs)
	return []audio.PCM16(samples)[:read], config, duration, nil
}

func cleanupAudio() {
	audioCancel()
	audioDevice.DeleteSources(int32(len(audioSources)), &audioSources[0])
	for _, s := range Files.sfx {
		s.Delete()
	}
	for _, s := range Files.music {
		s.Delete()
	}
	if audioDevice != nil {
		audioDevice.Close()
	}
}

func init() {
	al.SetErrorHandler(func(err error) {
		log.Println("[audio]", err)
	})
}

var _ Sound = (*Music)(nil)
var _ Sound = (*SFX)(nil)
