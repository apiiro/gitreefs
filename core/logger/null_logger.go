package logger

import (
	"io"
	"log"
)

type nullWriter struct {
}

func (writer *nullWriter) Write(buff []byte) (int, error) {
	return len(buff), nil
}

var _ io.Writer = &nullWriter{}

func newNullLogger() *log.Logger {
	return log.New(&nullWriter{}, "", 0)
}
