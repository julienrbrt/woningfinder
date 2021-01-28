package util

import (
	"context"
	"io"
)

type cancelFuncReadCloser struct {
	readCloser io.ReadCloser
	cancelFunc context.CancelFunc
}

func NewCancelFuncReadCloser(readCloser io.ReadCloser, cancelFunc context.CancelFunc) io.ReadCloser {
	return &cancelFuncReadCloser{
		readCloser: readCloser,
		cancelFunc: cancelFunc,
	}
}

func (c *cancelFuncReadCloser) Read(p []byte) (int, error) {
	return c.readCloser.Read(p)
}

func (c *cancelFuncReadCloser) Close() error {
	err := c.readCloser.Close()
	c.cancelFunc()
	return err
}
