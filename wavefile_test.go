package wavefile_test

import (
	"math"
	"testing"

	"github.com/domsim1/wavefile"
)

const sampleRate = 44100
const fadeRate = 8000

func TestSaveFile(t *testing.T) {
	waveFile := wavefile.NewWaveFile()
	err := waveFile.Create("test.wav")
	if err != nil {
		t.Error(err)
	}
	defer waveFile.Close()

	volume := 32000
	fadeInRate := volume / fadeRate

	var waveform [sampleRate]uint16

	frequency := 130.81
	v := 0
	for i := 0; i < sampleRate; i++ {
		if v < volume && i < sampleRate-fadeRate {
			v += fadeInRate
		}
		if i > sampleRate-fadeRate {
			v -= fadeInRate
		}
		t := float64(i) / float64(sampleRate)
		waveform[i] = uint16(float64(v) * math.Sin((frequency * t * 2 * math.Pi)))
	}
	err = waveFile.Write(waveform)
	if err != nil {
		t.Error(err)
	}

	frequency = 164.81
	v = 0
	for i := 0; i < sampleRate; i++ {
		if v < volume && i < sampleRate-fadeRate {
			v += fadeInRate
		}
		if i > sampleRate-fadeRate {
			v -= fadeInRate
		}
		t := float64(i) / float64(sampleRate)
		waveform[i] = uint16(float64(v) * math.Sin((frequency * t * 2 * math.Pi)))
	}
	err = waveFile.Write(waveform)
	if err != nil {
		t.Error(err)
	}

	frequency = 196.00
	v = 0
	for i := 0; i < sampleRate; i++ {
		if v < volume && i < sampleRate-fadeRate {
			v += fadeInRate
		}
		if i > sampleRate-fadeRate {
			v -= fadeInRate
		}
		t := float64(i) / float64(sampleRate)
		waveform[i] = uint16(float64(v) * math.Sin((frequency * t * 2 * math.Pi)))
	}
	err = waveFile.Write(waveform)
	if err != nil {
		t.Error(err)
	}

}
