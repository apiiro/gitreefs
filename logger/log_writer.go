package logger

import (
	"encoding/json"
	"io"
	"time"
)

type logEntry struct {
	Name             string `json:"name,omitempty"`
	LevelName        string `json:"levelname,omitempty"`
	Severity         string `json:"severity,omitempty"`
	Message          string `json:"message,omitempty"`
	TimestampSeconds int64  `json:"timestampSeconds,omitempty"`
	TimestampNanos   int    `json:"timestampNanos,omitempty"`
}

type jsonWriter struct {
	writers []io.Writer
	level   string
}

func (f *jsonWriter) Write(message []byte) (bytesWritten int, err error) {
	now := time.Now()

	entry := logEntry{
		Name:             "root",
		LevelName:        f.level,
		Severity:         f.level,
		Message:          string(message),
		TimestampSeconds: now.Unix(),
		TimestampNanos:   now.Nanosecond(),
	}

	var buf []byte
	buf, err = json.Marshal(entry)
	if err != nil {
		return
	}

	buf = append(buf, '\n')
	for _, writer := range f.writers {
		_, err = writer.Write(buf)
	}

	bytesWritten = len(message)
	return
}
