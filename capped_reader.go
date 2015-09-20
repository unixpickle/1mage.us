package main

import (
	"errors"
	"io"
)

var ErrCapExceeded = errors.New("size cap reached")

// CappedReader is an io.Reader which fails with an error when too much data is read.
type CappedReader struct {
	reader    io.Reader
	remaining int64
}

// NewCappedReader creates a CappedReader which uses the global database's size cap.
func NewCappedReader(input io.Reader) *CappedReader {
	return &CappedReader{input, GlobalDb.Config().MaxFileSize}
}

func (c *CappedReader) Read(p []byte) (int, error) {
	count, err := c.reader.Read(p)
	if err != nil && err != io.EOF {
		return 0, err
	} else {
		if int64(count) > c.remaining {
			return 0, ErrCapExceeded
		}
		c.remaining -= int64(count)
		return count, err
	}
}
