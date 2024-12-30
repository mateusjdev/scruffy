package common

import "mateusjdev/scruffy/cmd/clog"

func CheckIfError(err error) {
	if err != nil {
		clog.Errorf("%s\n", err)
		clog.ExitBecause(clog.CODE_ERROR)
	}
}
