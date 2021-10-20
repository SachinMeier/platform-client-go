// Package log is a simple logger interface
package log

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	// Production defines logging for production environments.
	Production = "production"
	// Development defines logging for development environments.
	Development = "development"
)

var env string

func Init(environment string) {
	switch environment {
	case Production:
		env = Production
	case Development:
		fallthrough
	default:
		env = Development
	}
}

type logWriter struct{}

func (f logWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}
	log.Printf("%s:%d :%s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}

var (
	logErr   = log.New(logWriter{}, "[ERR]: ", 0)
	logWarn  = log.New(logWriter{}, "[WARN]: ", 0)
	logInfo  = log.New(logWriter{}, "[INFO]: ", 0)
	logDebug = log.New(logWriter{}, "[DEBUG]: ", 0)
)

// Debug logs a message with debug log level.
func Debug(msg string) {
	if env == Development {
		logDebug.Println(msg)
	}
}

// Debugf logs a formatted message with debug log level.
func Debugf(msg string, args ...interface{}) {
	if env == Development {
		logDebug.Println(fmt.Sprintf(msg, args...))
	}
}

// Info logs a message with info log level.
func Info(msg string) {
	logInfo.Println(msg)
}

// Infof logs a formatted message with info log level.
func Infof(msg string, args ...interface{}) {
	logInfo.Println(fmt.Sprintf(msg, args...))
}

// Warn logs a message with warn log level.
func Warn(msg string) {
	logWarn.Println(msg)
}

// Warnf logs a formatted message with warn log level.
func Warnf(msg string, args ...interface{}) {
	logWarn.Println(fmt.Sprintf(msg, args...))
}

// Error logs a message with error log level.
func Error(msg interface{}) {
	logErr.Println(msg)
}

// Errorf logs a formatted message with error log level.
func Errorf(msg string, args ...interface{}) {
	logErr.Println(fmt.Sprintf(msg, args...))
}

// Fatal logs a message with fatal log level.
func Fatal(msg string) {
	log.Fatalln(msg)
}

// Fatalf logs a formatted message with fatal log level.
func Fatalf(msg string, args ...interface{}) {
	log.Fatalln(fmt.Sprintf(msg, args...))
}

// Query logs a query statement with a debugf log level.
func Query(query string, args ...interface{}) {
	if env == Development {
		// Replace `$1` OR `?` with `%v`
		reg, _ := regexp.Compile(`\$[0-9]|\?`)
		processedString := reg.ReplaceAllString(query, "%v")
		Debugf(processedString, args)
	}
}
