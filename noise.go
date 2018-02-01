// Package noise makes noise from a keyboards input.
// Based off of https://github.com/go-audio/generator/blob/master/examples/realtime/main.go
package noise

import (
	"log"

	"github.com/go-audio/audio"
	"github.com/go-audio/generator"
	"github.com/gordonklaus/portaudio"
	"golang.org/x/crypto/ssh/terminal"
)

// Noise holds a channel for read characters from stdin that
// are then used to control sound.
type Noise struct {
	buf               *audio.FloatBuffer
	bufferSize        int
	currentNote       float64
	currentVol        float64
	gainChan          chan float64
	keyChan           chan rune
	osc               *generator.Osc
	output            []float32
	stopGeneratorChan chan struct{}
	stopReadChan      chan struct{}
	stream            *portaudio.Stream
	waveTypeChan      chan generator.WaveType
}

// New initalizes a new Noise, reading stdin in the background
// and generating a signal based on input.
func New() (*Noise, error) {
	bufferSize := 1024
	buf := &audio.FloatBuffer{
		Data:   make([]float64, bufferSize),
		Format: audio.FormatMono44100,
	}
	currentNote := 440.0
	osc := generator.NewOsc(generator.WaveSine, currentNote, buf.Format.SampleRate)
	osc.Amplitude = 0.5

	n := &Noise{
		buf:               buf,
		bufferSize:        bufferSize,
		currentNote:       440.0,
		currentVol:        osc.Amplitude,
		gainChan:          make(chan float64),
		keyChan:           make(chan rune),
		osc:               osc,
		output:            make([]float32, bufferSize),
		stopGeneratorChan: make(chan struct{}),
		stopReadChan:      make(chan struct{}),
		waveTypeChan:      make(chan generator.WaveType),
	}
	return n, nil
}

// Run sets up the terminal and initialize the audio stream.
func (n *Noise) Run() error {
	// setup the terminal in raw mode
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)

	// run the stdin reader in the background
	go func() {
		if err := n.read(); err != nil {
			log.Fatal(err)
		}
	}()

	// initialize portaudio
	portaudio.Initialize()
	defer portaudio.Terminate()

	// initialize and start the audio stream
	n.stream, err = portaudio.OpenDefaultStream(0, 1, 44100, len(n.output), &n.output)
	if err != nil {
		return err
	}
	defer n.stream.Close()
	if err := n.stream.Start(); err != nil {
		return err
	}
	defer n.stream.Stop()

	// run the audio generator in the background
	go n.generator()

	// wait for esc to quit
	for {
		select {
		case <-n.stopReadChan:
			n.stopGeneratorChan <- struct{}{}
			return nil
		}
	}
}
