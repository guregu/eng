package engi

import (
	"log"
	"os"
	"time"
	"unsafe"

	"azul3d.org/audio.v1"
	"azul3d.org/native/al.v1"

	_ "azul3d.org/audio/flac.dev"
	_ "azul3d.org/audio/wav.v1"
)

var (
	audioDevice  *al.Device
	audioSources []uint32
	muted        bool
)

const (
	audioSourcesCount = 31
	buffersPerStream  = 3
)

type Stream struct {
	source  uint32
	buffers []uint32
	looping bool
	done    chan struct{}

	decoder audio.Decoder
	file    *os.File
}

func (s *Stream) Play() {
	s.looping = false
	s.play()
}

func (s *Stream) Loop() {
	s.looping = true
	s.play()
}

func (s *Stream) play() {
	s.done = make(chan struct{})
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

func (s *Stream) reset() {
	// flac pkg can't seek yet so ghetto seek to 0
	s.file.Seek(0, os.SEEK_SET)
	var err error
	s.decoder, _, err = audio.NewDecoder(s.file)
	if err != nil {
		panic(err)
	}
}

func (s *Stream) run() {
	defer func() {
		s.done = nil
	}()
	for {
		select {
		case <-s.done:
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

func (s *Stream) unqueue() []uint32 {
	var processed int32

	audioDevice.GetSourcei(s.source, al.BUFFERS_PROCESSED, &processed)
	if processed == 0 {
		return nil
	}
	available := make([]uint32, processed)
	audioDevice.SourceUnqueueBuffers(s.source, available)
	return available
}

func (s *Stream) fill(buffer uint32) error {
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

func (s *Stream) Delete() {
	if s.done != nil {
		close(s.done)
	}
	s.file.Close()

	audioDevice.DeleteBuffers(buffersPerStream, &s.buffers[0])
}

func (s *Stream) Playing() bool {
	if s.source == 0 {
		return false
	}
	var state int32
	audioDevice.GetSourcei(s.source, al.SOURCE_STATE, &state)
	return state == al.PLAYING
}

func (s *Stream) Stop() {
	if s.done != nil {
		close(s.done)
	}
	audioDevice.SourceStop(s.source)
	s.unqueue()
	s.source = 0
}

func loadStream(r Resource) (*Stream, error) {
	// func readSoundFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {
	file, err := os.Open(r.url)
	if err != nil {
		return nil, err
	}

	s := Stream{
		buffers: make([]uint32, buffersPerStream),
		file:    file,
	}

	audioDevice.GenBuffers(buffersPerStream, &s.buffers[0])

	return &s, nil
}

type Sound struct {
	source     uint32
	buffer     uint32
	duration   time.Duration // seconds
	sampleRate int
	looping    bool
}

func (s *Sound) bind() {
	s.source = nextAvailableSource()
	audioDevice.Sourcei(s.source, al.BUFFER, int32(s.buffer))
	if s.looping {
		audioDevice.Sourcei(s.source, al.LOOPING, al.TRUE)
	} else {
		audioDevice.Sourcei(s.source, al.LOOPING, al.FALSE)
	}
}

func (s *Sound) Play() {
	s.looping = false
	s.bind()
	audioDevice.SourcePlay(s.source)
}

func (s *Sound) Loop() {
	s.looping = true
	s.bind()
	audioDevice.SourcePlay(s.source)
}

func (s *Sound) Stop() {
	audioDevice.SourceStop(s.source)
	s.unqueue()
}

func (s *Sound) unqueue() []uint32 {
	var processed int32
	audioDevice.GetSourcei(s.source, al.BUFFERS_PROCESSED, &processed)
	if processed == 0 {
		return nil
	}
	available := make([]uint32, processed)
	audioDevice.SourceUnqueueBuffers(s.source, available)
	return available
}

func (s *Sound) Delete() {
	audioDevice.DeleteBuffers(1, &s.buffer)
}

func (s *Sound) Playing() bool {
	var state int32
	audioDevice.GetSourcei(s.source, al.SOURCE_STATE, &state)
	return state == al.PLAYING
}

func (s *Sound) Duration() time.Duration {
	return s.duration
}

func ToggleMute() {
	muted = !muted
	gain := float32(1)
	if muted {
		gain = 0
	}

	for _, src := range audioSources {
		audioDevice.Sourcef(src, al.GAIN, gain)
	}
}

func Muted() bool {
	return muted
}

func setupAudio() {
	var err error
	audioDevice, err = al.OpenDevice("", nil)
	fatalErr(err)

	audioSources = make([]uint32, audioSourcesCount)
	audioDevice.GenSources(audioSourcesCount, &audioSources[0])
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

func loadSound(r Resource) (*Sound, error) {
	samples, config, duration, err := readSoundFile(r.url)
	if err != nil {
		return nil, err
	}

	s := Sound{
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

func readSoundFile(filename string) (samples []audio.PCM16, config audio.Config, duration time.Duration, err error) {
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
	audioDevice.DeleteSources(int32(len(audioSources)), &audioSources[0])
	for _, s := range Files.sounds {
		if s.Playing() {
			s.Stop()
		}
		s.Delete()
	}
	for _, s := range Files.streams {
		if s.Playing() {
			s.Stop()
		}
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
