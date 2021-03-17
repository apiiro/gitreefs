package logger

import (
	"fmt"
	"io"
	"time"
)

type logWriter struct {
	io.Writer
	writers []io.Writer
	level   string
}

func (lw *logWriter) Write(message []byte) (bytesWritten int, err error) {
	endOfLine := "\n"
	if message[len(message)-1] == '\n' {
		endOfLine = ""
	}
	formattedMessage := []byte(fmt.Sprintf(
		"%v [%v] %v: %s%s",
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		lw.level,
		appVersion,
		message,
		endOfLine,
	))

	for _, writer := range lw.writers {
		if writer != nil {
			_, err = writer.Write(formattedMessage)
		}
	}

	bytesWritten = len(message)
	return
}
