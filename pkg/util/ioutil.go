package util

import (
	"io"
	"io/ioutil"
)

func ReadAllAndClose(r io.ReadCloser) ([]byte, error) {
	defer func() {
		_ = r.Close()
	}()

	return ioutil.ReadAll(r)
}
