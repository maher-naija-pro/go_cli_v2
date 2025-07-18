package logger

import (
    "io"
    "log"
    "os"
    "sync"
)

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
)

var (
    logger   *log.Logger
    logLevel LogLevel = INFO
    once     sync.Once
)

// InitLogger initializes the logger with the given output and log level.
func InitLogger(out io.Writer, level LogLevel) {
    once.Do(func() {
        logger = log.New(out, "", log.LstdFlags)
        logLevel = level
    })
}

// SetLogLevel allows changing the log level at runtime.
func SetLogLevel(level LogLevel) {
    logLevel = level
}

// SetLogOutput allows changing the log output at runtime.
func SetLogOutput(out io.Writer) {
    logger.SetOutput(out)
}

func logf(level LogLevel, format string, v ...interface{}) {
    if logger == nil {
        InitLogger(os.Stderr, INFO)
    }
    if level >= logLevel {
        logger.Printf(format, v...)
    }
}

func Debugf(format string, v ...interface{}) { logf(DEBUG, "[DEBUG] "+format, v...) }
func Infof(format string, v ...interface{})  { logf(INFO, "[INFO] "+format, v...) }
func Warnf(format string, v ...interface{})  { logf(WARN, "[WARN] "+format, v...) }
func Errorf(format string, v ...interface{}) { logf(ERROR, "[ERROR] "+format, v...) } 