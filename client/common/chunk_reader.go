package common

import (
	"bufio"
	"os"
)

type ChunkReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	chunkSize int
}

func newChunkReader(path string, chuckSize int) (*ChunkReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	ChunkReader := &ChunkReader{file: file, scanner: scanner, chunkSize: chuckSize}
	return ChunkReader, nil
}

func (self *ChunkReader) read() ([]string, error) {
	lines := make([]string, 0)

	for i := 0; i < self.chunkSize && self.scanner.Scan(); i++ {
		line := self.scanner.Text()
		lines = append(lines, line)
	}

	return lines, self.scanner.Err()
}

func (self *ChunkReader) close() {
	self.file.Close()
}
