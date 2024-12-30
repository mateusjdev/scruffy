package clog

import "os"

type exitCode uint8

const (
	ErrCodeGeneric exitCode = iota + 1
	ErrUserGeneric
)

func ExitBecause(reason exitCode) {
	os.Exit(int(reason))
}

func CheckIfError(err error) {
	if err != nil {
		Errorf("%s\n", err)
		ExitBecause(ErrCodeGeneric)
	}
}
