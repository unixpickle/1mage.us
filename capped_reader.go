package main

import (
	"errors"
	"io"
)

var ErrCapExceeded = errors.New("size cap reached")

// CappedReader is an io.Reader which fails with an error when too much data is read.
type CappedReader struct {
	reader    io.Reader
	remaining int
}

// NewCappedReader creates a CappedReader which uses the global database's size cap.
func NewCappedReader(input io.Reader) *CappedReader {
	GlobalDatabase.RLock()
	maxSize := GlobalDatabase.Config.MaxFileSize
	GlobalDatabase.RUnlock()
	return &CappedReader{input, maxSize}
}

func (c *CappedReader) Read(p []byte) (int, error) {
	count, err := c.reader.Read(p)
	if err != nil && err != io.EOF {
		return 0, err
	} else {
		if count > c.remaining {
			return 0, ErrCapExceeded
		}
		c.remaining -= count
		return count, err
	}
}
