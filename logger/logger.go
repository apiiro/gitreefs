package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogLevel uint32

const (
	_             LogLevel = iota
	LogLevelDebug          = iota
	LogLevelInfo           = iota
	LogLevelError          = iota
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	File        *os.File
	GlobalLevel LogLevel
)

func InitLoggers(filePathFormat string, level string) error {
	var filePath = fmt.Sprintf(filePathFormat, time.Now().UTC().Format("yyyy-MMM-dd"))
	var filePermissions os.FileMode = 0644
	err := os.Mkdir(filepath.Dir(filePath), filePermissions)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		filePermissions,
	)
	if err != nil {
		return err
	}

	File = file
	GlobalLevel = stringToLevel(level)

	var flag = log.Ldate | log.Ltime | log.Lmicroseconds
	DebugLogger = log.New(createLogger(LogLevelDebug), "DEBUG: ", flag)
	InfoLogger = log.New(createLogger(LogLevelInfo), "INFO: ", flag)
	ErrorLogger = log.New(createLogger(LogLevelError), "ERROR: ", flag)

	return nil
}

func levelToString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelError:
		return "ERROR"
	default:
		return "unknown"
	}
}

func stringToLevel(level string) LogLevel {
	switch level {
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelDebug
	case "ERROR":
		return LogLevelDebug
	default:
		return 0
	}
}

func createLogger(level LogLevel) io.Writer {

	var consoleWriter io.Writer
	switch level {
	case LogLevelError:
		consoleWriter = os.Stderr
	default:
		consoleWriter = os.Stdout
	}

	return &jsonWriter{
		level:   levelToString(level),
		writers: []io.Writer{consoleWriter, File},
	}
}

func CloseLoggers() {
	if File != nil {
		File.Close()
		File = nil
	}
}

func Debug(v ...interface{}) {
	if GlobalLevel > LogLevelDebug {
		return
	}
	DebugLogger.Println(v...)
}

func Info(v ...interface{}) {
	if GlobalLevel > LogLevelInfo {
		return
	}
	InfoLogger.Println(v...)
}

func Error(v ...interface{}) {
	if GlobalLevel > LogLevelError {
		return
	}
	ErrorLogger.Println(v...)
}
