package wavefile

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type waveFileHeader struct {
	chunkID       [4]byte
	chunkSize     uint32
	format        [4]byte
	subchunk1ID   [4]byte
	subchunk1Size uint32
	audioFormat   uint16
	numChannels   uint16
	sampleRate    uint32
	byteRate      uint32
	blockAlign    uint16
	bitsPerSample uint16
	subchunk2ID   [4]byte
	subchunk2Size uint32
}

type waveFile struct {
	header *waveFileHeader
	file   *os.File
}

type IWaveFile interface {
	Create(fileName string) error
	Close() error
	Write(data interface{}) error
}

type WaveFileSetting struct {
	SampleRate       uint32
	BitsPerSample    uint16
	NumberOfChannels uint16
}

func newWaveFileHeader(settings *WaveFileSetting) *waveFileHeader {
	w := &waveFileHeader{
		chunkID:       [4]byte{'R', 'I', 'F', 'F'},
		chunkSize:     0,
		format:        [4]byte{'W', 'A', 'V', 'E'},
		subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		subchunk1Size: 16,
		audioFormat:   1,
		numChannels:   settings.NumberOfChannels,
		sampleRate:    settings.SampleRate,
		bitsPerSample: settings.BitsPerSample,
		subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		subchunk2Size: 0,
	}
	updateComputedFields(w)
	return w
}

func setByteRate(w *waveFileHeader) {
	w.byteRate = (w.sampleRate * uint32(w.bitsPerSample) * uint32(w.numChannels)) / 8
}

func setBlockAlign(w *waveFileHeader) {
	w.blockAlign = (w.bitsPerSample * w.numChannels) / 8
}

func updateComputedFields(w *waveFileHeader) {
	setByteRate(w)
	setBlockAlign(w)
}

func NewWaveFile() IWaveFile {
	w := &waveFile{
		header: newWaveFileHeader(&WaveFileSetting{
			SampleRate:       44100,
			BitsPerSample:    16,
			NumberOfChannels: 1,
		}),
	}
	return w
}

func NewWaveFileWithSettings(settings *WaveFileSetting) IWaveFile {
	w := &waveFile{
		header: newWaveFileHeader(settings),
	}
	return w
}

func (waveFile *waveFile) Create(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	waveFile.file = file
	err = waveFile.Write(*waveFile.header)
	if err != nil {
		return err
	}

	return nil
}

func (waveFile *waveFile) Close() error {
	fileLength, err := waveFile.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	dataLength := uint32(fileLength - 44)
	waveFile.file.Seek(40, io.SeekStart)
	waveFile.Write(dataLength)
	riffLength := uint32(fileLength - 8)
	waveFile.file.Seek(4, io.SeekStart)
	waveFile.Write(riffLength)
	err = waveFile.file.Close()
	waveFile.file = nil
	return err
}

func (waveFile *waveFile) Write(data interface{}) error {
	if waveFile.file == nil {
		return fmt.Errorf("could not find open file")
	}
	return binary.Write(waveFile.file, binary.LittleEndian, data)
}
