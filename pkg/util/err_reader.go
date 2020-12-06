package util

import "io"

type ErrReader interface {
	io.Reader
	Error() error
}

func NewErrReader(err error) ErrReader {
	return &errReader{err: err}
}

type errReader struct {
	err error
}

func (r *errReader) Read([]byte) (int, error) {
	return 0, r.err
}

func (r *errReader) Error() error {
	return r.err
}
