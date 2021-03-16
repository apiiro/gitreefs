package common

type Options interface {
	LogFile() string
	LogLevel() string
}
