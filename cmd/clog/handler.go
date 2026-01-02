package clog

import "os"

type exitCode uint8

const (
	ErrCodeGeneric exitCode = iota + 1
	ErrUserGeneric
	ErrUserInput
)

func PanicReturning(err error, reason exitCode) {
	os.Exit(int(reason))
}

func PanicIf(err error) {
	if err != nil {
		Errorf("%s\n", err)
		os.Exit(int(ErrCodeGeneric))
	}
}
