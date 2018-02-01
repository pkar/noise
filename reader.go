package noise

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/go-audio/generator"
)

// Reads character by character from stdin. esc will exit reading,
// +- control the volume, and all other characters change the frequency
// of the oscillator.
func (n *Noise) read() error {
	fmt.Printf("Commands:\n\resc: to quit\n\r" +
		"+/-: for volume control\n\r" +
		"0: WaveSine\n\r" +
		"1: WaveTriangle\n\r" +
		"2: WaveSaw\n\r3: WaveSqr\n\r" +
		"any other character to change the frequency\n\r")

	// read character by character.
	r := bufio.NewReader(os.Stdin)
	for {
		k, _, err := r.ReadRune()
		if err != nil {
			return err
		}

		switch k {
		case '\x1b':
			n.stopReadChan <- struct{}{}
			return nil
		case '+':
			n.gainChan <- 0.10
		case '-':
			n.gainChan <- -0.10
		case '0', '1', '2', '3':
			v, _ := strconv.Atoi(string(k))
			n.waveTypeChan <- generator.WaveType(v)
		default:
			n.keyChan <- k
		}
	}
}
