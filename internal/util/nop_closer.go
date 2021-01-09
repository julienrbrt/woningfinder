package util

import (
	"errors"
	"io"
)

type NopCloser interface {
	io.ReadCloser
	Closed() bool
}

type nopCloser struct {
	reader io.Reader
	closed bool
}

func NewNopCloser(r io.Reader) NopCloser {
	return &nopCloser{
		reader: r,
		closed: false,
	}
}

func (c *nopCloser) Read(p []byte) (int, error) {
	if c.closed {
		return 0, errors.New("trying to read from closed ReadCloser")
	}
	return c.reader.Read(p)
}

func (c *nopCloser) Close() error {
	c.closed = true
	c.reader = nil
	return nil
}

func (c *nopCloser) Closed() bool {
	return c.closed
}
