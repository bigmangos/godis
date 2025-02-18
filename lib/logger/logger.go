package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Settings stores config for logger
type Settings struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time-format"`
}

var (
	level              = ERROR
	defaultPrefix      = ""
	defaultCallerDepth = 2
	logger             *log.Logger
	mu                 sync.Mutex
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

type logLevel int

// log levels
const (
	CLOSE logLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const flags = log.LstdFlags

func init() {
	logger = log.New(os.Stdout, defaultPrefix, flags)
}

// Setup initializes logger
func Setup(settings *Settings) {
	var err error
	dir := settings.Path
	fileName := fmt.Sprintf("%s-%s.%s",
		settings.Name,
		time.Now().Format(settings.TimeFormat),
		settings.Ext)

	logFile, err := mustOpen(fileName, dir)
	if err != nil {
		log.Fatalf("logging.Setup err: %s", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, defaultPrefix, flags)
}

func setPrefix(level logLevel) {
	_, file, line, ok := runtime.Caller(defaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d] ", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s] ", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

// Debug prints debug log
func Debug(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if level > DEBUG {
		return
	}
	setPrefix(DEBUG)
	logger.Printf(format, v...)
}

// Info prints normal log
func Info(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if level > INFO {
		return
	}
	setPrefix(INFO)
	logger.Printf(format, v)
}

// Warn prints warning log
func Warn(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if level > WARNING {
		return
	}
	setPrefix(WARNING)
	logger.Printf(format, v...)
}

// Error prints error log
func Error(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if level > ERROR {
		return
	}
	setPrefix(ERROR)
	logger.Printf(format, v...)
}

// Fatal prints error log then stop the program
func Fatal(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	setPrefix(FATAL)
	logger.Printf(format, v...)
}
