package logger

import (
    "fmt"
    "io"
    "log"
    "os"
    "runtime"
    "strings"
    "sync"
    "time"
)

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

var (
    logger   *log.Logger
    logLevel LogLevel = INFO
    once     sync.Once
    mu       sync.Mutex
)

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
    switch l {
    case DEBUG:
        return "DEBUG"
    case INFO:
        return "INFO"
    case WARN:
        return "WARN"
    case ERROR:
        return "ERROR"
    case FATAL:
        return "FATAL"
    default:
        return "UNKNOWN"
    }
}

// InitLogger initializes the logger with the given output and log level.
func InitLogger(out io.Writer, level LogLevel) {
    once.Do(func() {
        logger = log.New(out, "", 0)
        logLevel = level
    })
}

// SetLogLevel allows changing the log level at runtime.
func SetLogLevel(level LogLevel) {
    mu.Lock()
    defer mu.Unlock()
    logLevel = level
}

// SetLogOutput allows changing the log output at runtime.
func SetLogOutput(out io.Writer) {
    mu.Lock()
    defer mu.Unlock()
    if logger == nil {
        InitLogger(out, logLevel)
    } else {
        logger.SetOutput(out)
    }
}

// logf is the internal logging function with enhanced formatting.
func logf(level LogLevel, format string, v ...interface{}) {
    mu.Lock()
    defer mu.Unlock()
    if logger == nil {
        InitLogger(os.Stderr, INFO)
    }
    if level < logLevel {
        return
    }
    now := time.Now().Format("2006-01-02 15:04:05")
    pc, file, line, ok := runtime.Caller(3)
    var caller string
    if ok {
        fn := runtime.FuncForPC(pc)
        if fn != nil {
            caller = fmt.Sprintf("%s:%d %s", shortFile(file), line, fn.Name())
        } else {
            caller = fmt.Sprintf("%s:%d", shortFile(file), line)
        }
    } else {
        caller = "unknown"
    }
    prefix := fmt.Sprintf("[%s] [%s] [%s]", now, level.String(), caller)
    msg := fmt.Sprintf(format, v...)
    logger.Printf("%s %s", prefix, msg)
    if level == FATAL {
        os.Exit(1)
    }
}

// shortFile returns the last part of the file path.
func shortFile(path string) string {
    for i := len(path) - 1; i >= 0; i-- {
        if path[i] == '/' || path[i] == '\\' {
            return path[i+1:]
        }
    }
    return path
}

func Debugf(format string, v ...interface{}) { logf(DEBUG, format, v...) }
func Infof(format string, v ...interface{})  { logf(INFO, format, v...) }
func Warnf(format string, v ...interface{})  { logf(WARN, format, v...) }
func Errorf(format string, v ...interface{}) { logf(ERROR, format, v...) }
func Fatalf(format string, v ...interface{}) { logf(FATAL, format, v...) }

func ParseLogLevel(s string) (LogLevel, error) {
    switch strings.ToUpper(s) {
    case "DEBUG":
        return DEBUG, nil
    case "INFO":
        return INFO, nil
    case "WARN", "WARNING":
        return WARN, nil
    case "ERROR":
        return ERROR, nil
    case "FATAL":
        return FATAL, nil
    default:
        return INFO, fmt.Errorf("unknown log level: %s", s)
    }
}