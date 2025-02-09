package lib

import (
	"bytes"
	"io"
)

func GetSize(stream io.Reader) int {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Len()
}
