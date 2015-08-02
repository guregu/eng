package engi

import (
	"log"
	"os"
	"unsafe"

	"azul3d.org/audio.v1"
	"azul3d.org/native/al.v1"

	_ "azul3d.org/audio/flac.dev"
	_ "azul3d.org/audio/wav.v1"
)

var audioDevice *al.Device

type Sound struct {
	source     uint32
	buffer     uint32
	duration   float64 // seconds
	sampleRate int
	looping    bool
}

func (s *Sound) Play() {
	if s.looping {
		audioDevice.Sourcei(s.source, al.LOOPING, al.FALSE)
	}
	audioDevice.SourcePlay(s.source)
}

func (s *Sound) Loop() {
	if !s.looping {
		audioDevice.Sourcei(s.source, al.LOOPING, al.TRUE)
	}
	audioDevice.SourcePlay(s.source)
}

func (s *Sound) Stop() {
	audioDevice.SourceStop(s.source)
}

func (s *Sound) Delete() {
	audioDevice.DeleteSources(1, &s.source)
	audioDevice.DeleteBuffers(1, &s.buffer)
}

func (s *Sound) Playing() bool {
	var state int32
	audioDevice.GetSourcei(s.source, al.SOURCE_STATE, &state)
	return state == al.PLAYING || state == al.LOOPING
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
	audioDevice.GenSources(1, &s.source)
	audioDevice.GenBuffers(1, &s.buffer)
	if config.Channels == 1 {
		audioDevice.BufferData(s.buffer, al.FORMAT_MONO16, unsafe.Pointer(&samples[0]), int32(int(unsafe.Sizeof(samples[0]))*len(samples)), int32(config.SampleRate))
	} else {
		audioDevice.BufferData(s.buffer, al.FORMAT_STEREO16, unsafe.Pointer(&samples[0]), int32(int(unsafe.Sizeof(samples[0]))*len(samples)), int32(config.SampleRate))
	}
	audioDevice.Sourcei(s.source, al.BUFFER, int32(s.buffer))
	return &s, err
}

func readSoundFile(filename string) (samples []audio.PCM16, config audio.Config, duration float64, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, audio.Config{}, 0, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, audio.Config{}, 0, err
	}

	decoder, _, err := audio.NewDecoder(file)
	if err != nil {
		return nil, audio.Config{}, 0, err
	}

	config = decoder.Config()
	// Guess duration (mostly accurate for WAVs)
	duration = float64(fi.Size()) / float64(config.SampleRate*config.Channels*16/8)

	// Create a buffer that can hold 3 second of audio samples
	// bufSize := int(duration * float64(config.SampleRate*config.Channels)) // undersized for flac files
	// Convert everything to 16-bit samples
	// fmt.Println("bufsize", 1024*1024*5)
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
	err = nil

	duration = 1 / float64(config.SampleRate) * float64(read)
	return []audio.PCM16(samples)[:read], config, float64(duration), err
}

func init() {
	al.SetErrorHandler(func(err error) {
		log.Println("[audio]", err)
	})
}
