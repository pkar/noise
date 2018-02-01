package noise

import (
	"fmt"
	"math"

	"github.com/go-audio/transforms"
)

func (n *Noise) generator() {
	for {
		select {
		case g := <-n.gainChan:
			n.currentVol += g
			if n.currentVol < 0.1 {
				n.currentVol = 0
			}
			if n.currentVol > 6 {
				n.currentVol = 6
			}
			fmt.Printf("new vol %f.2\n\r", n.currentVol)
		case k := <-n.keyChan:
			v := float64(math.Abs(float64(int(k - 100))))
			newNote := 440.0 * math.Pow(2, (v)/12.0)
			if newNote > 22000 {
				fmt.Printf("dropping note change above 22000 for %q %.2f Hz\n\r", k, n.currentNote)
				continue
			}
			if n.currentNote != newNote {
				fmt.Printf("switching oscillator to %.2f Hz\n\r", n.currentNote)
				n.currentNote = newNote
				n.osc.SetFreq(n.currentNote)
			}
			// populate the out buffer
			if err := n.osc.Fill(n.buf); err != nil {
				fmt.Printf("error filling up the buffer\n\r")
			}
			transforms.Gain(n.buf, n.currentVol)

			f64ToF32Copy(n.output, n.buf.Data)

			// write to the stream
			if err := n.stream.Write(); err != nil {
				//fmt.Printf("error writing to stream: %v\n\r", err)
			}
		case w := <-n.waveTypeChan:
			n.osc.Shape = w
		case <-n.stopGeneratorChan:
			return
		}
	}
}
