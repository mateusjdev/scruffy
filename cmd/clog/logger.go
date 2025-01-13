package clog

import (
	"fmt"
	"os"
	"strings"
)

// TODO: use a custom lib?
// TODO: add a --no-color or --simple-text or --json

var logLoggerLevel logLevel

type logLevel uint8

const (
	LevelDebug logLevel = iota
	LevelInfo
	LevelWarning
	LevelError
)

const (
	printDebug       string = "\x1b[30;47m DEBUG \x1b[0m "
	printInfo        string = "\x1b[37;44m INFO \x1b[0m "
	printInfoSuccess string = "\x1b[37;42m INFO \x1b[0m "
	printInfoError   string = "\x1b[37;41m INFO \x1b[0m "
	printWarning     string = "\x1b[37;43m WARNING \x1b[0m "
	printError       string = "\x1b[37;41m ERROR \x1b[0m "
)

// Ensure space between logs:
// Dumb, but if works, it works
func ensureNewLine(msg *string) {
	if strings.HasSuffix(*msg, "\n") {
		return
	}
	*msg += "\n"
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
	levelPrintfOut(LevelDebug, printDebug+msg, args...)
}

func Infof(msg string, args ...any) {
	levelPrintfOut(LevelInfo, printInfo+msg, args...)
}

func InfoSuccessf(msg string, args ...any) {
	levelPrintfOut(LevelInfo, printInfoSuccess+msg, args...)
}

func InfoErrorf(msg string, args ...any) {
	levelPrintfOut(LevelInfo, printInfoError+msg, args...)
}

func Warningf(msg string, args ...any) {
	levelPrintfOut(LevelWarning, printWarning+msg, args...)
}

func Errorf(msg string, args ...any) {
	levelPrintfErr(LevelError, printError+msg, args...)
}

func SetLogLevel(level logLevel) {
	logLoggerLevel = level
}

func init() {
	logLoggerLevel = LevelInfo
}
