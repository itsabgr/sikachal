package common

import (
	"io"
)

func Throw(e error) {
	if e != nil {
		panic(e)
	}
}
func Must[R any](r R, e error) R {
	Throw(e)
	return r
}
func Close(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}
