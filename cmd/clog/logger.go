package clog

import (
	"fmt"
	"os"
	"strings"
)

var logLoggerLevel logLevel

type logLevel uint8

const (
	LevelDebug logLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	// LevelCritical
)

type StringIcon string

const (
	printDebug       StringIcon = "\x1b[30;47m DEBUG \x1b[0m "
	printInfo        StringIcon = "\x1b[30;44m INFO \x1b[0m "
	printWarning     StringIcon = "\x1b[30;43m WARNING \x1b[0m "
	printError       StringIcon = "\x1b[30;41m ERROR \x1b[0m "
	PrintIconSuccess StringIcon = "\x1b[37;42m ✔️ \x1b[0m "
	PrintIconNothing StringIcon = "\x1b[37;44m ➖ \x1b[0m "
	PrintIconError   StringIcon = "\x1b[37;41m ❌ \x1b[0m "
)

// Ensure space between logs:
// Dumb, but if works, it works
func ensureNewLine(msg *string) {
	if strings.HasSuffix(*msg, "\n") {
		return
	}
	*msg = *msg + "\n"
}

func levelPrintfOut(level logLevel, msg string, a ...any) {
	if logLoggerLevel <= level {
		ensureNewLine(&msg)
		fmt.Fprintf(os.Stdout, msg, a...)
	}
}

func levelPrintfErr(level logLevel, msg string, a ...any) {
	if logLoggerLevel <= level {
		ensureNewLine(&msg)
		fmt.Fprintf(os.Stderr, msg, a...)
	}
}

func Debugf(msg string, args ...any) {
	msg = string(printDebug) + msg
	levelPrintfOut(LevelDebug, msg, args...)
}

func Infof(msg string, args ...any) {
	msg = string(printInfo) + msg
	levelPrintfOut(LevelInfo, msg, args...)
}

func InfoIconf(icon StringIcon, msg string, args ...any) {
	msg = string(printInfo) + string(icon) + msg
	levelPrintfOut(LevelInfo, msg, args...)
}

func Warningf(msg string, args ...any) {
	msg = string(printWarning) + msg
	levelPrintfOut(LevelWarning, msg, args...)
}

func Errorf(msg string, args ...any) {
	msg = string(printError) + msg
	levelPrintfErr(LevelError, msg, args...)
}

func SetLogLevel(level logLevel) {
	logLoggerLevel = level
}

func init() {
	logLoggerLevel = LevelInfo
}
