package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LogLevel uint32

const (
	_ LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelError
)

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fileHandler *os.File
	globalLevel LogLevel
	appVersion  = ""
)

func init() {
	initLoggers()
}

func initLoggers() {
	flag := log.Ldate | log.Ltime | log.Lmicroseconds
	debugLogger = log.New(createLogger(LogLevelDebug), "", flag)
	infoLogger = log.New(createLogger(LogLevelInfo), "", flag)
	errorLogger = log.New(createLogger(LogLevelError), "", flag)
}

func InitLoggers(filePathFormat string, level string, version string) error {
	var filePath = fmt.Sprintf(filePathFormat, time.Now().UTC().Format("2006-01-02"), os.Getpid())
	var filePermissions os.FileMode = 0777
	err := os.Mkdir(filepath.Dir(filePath), filePermissions)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}

	file, err := os.OpenFile(
		filePath,
		os.O_CREATE|os.O_APPEND|os.O_RDWR,
		filePermissions,
	)
	if err != nil {
		return err
	}

	fileHandler = file
	globalLevel = stringToLevel(level)
	appVersion = version

	initLoggers()

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
		return LogLevelInfo
	case "ERROR":
		return LogLevelError
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

	writers := []io.Writer{consoleWriter}
	if fileHandler != nil {
		writers = append(writers, fileHandler)
	}
	return &logWriter{
		level:   levelToString(level),
		writers: writers,
	}
}

func CloseLoggers() {
	if fileHandler != nil {
		fileHandler.Close()
		fileHandler = nil
	}
}

func Debug(format string, v ...interface{}) {
	if globalLevel > LogLevelDebug {
		return
	}
	debugLogger.Printf(format, v...)
}

func Info(format string, v ...interface{}) {
	if globalLevel > LogLevelInfo {
		return
	}
	infoLogger.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	if globalLevel > LogLevelError {
		return
	}
	errorLogger.Printf(format, v...)
}

func DebugLogger() *log.Logger {
	if globalLevel > LogLevelDebug {
		return nil
	}
	return debugLogger
}

func ErrorLogger() *log.Logger {
	if globalLevel > LogLevelError {
		return nil
	}
	return errorLogger
}
