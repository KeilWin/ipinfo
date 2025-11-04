package utils

import (
	"log/slog"
	"os"
)

type ExitCodeType int

const (
	ExitSuccess ExitCodeType = 0
	ExitError   ExitCodeType = 1
)

func CheckLoadConfigError(err error, name string, component string) bool {
	if err != nil {
		slog.Error("can't load config variable", "component", component, "name", name, "error", err)
		return true
	}
	return false
}

func CheckAppFatalError(err error) {
	if err != nil {
		slog.Error("app fatal error", "error", err)
		os.Exit(int(ExitError))
	}
}

func CheckAppFatalManyErrors(errs []error, comment string) {
	if len(errs) != 0 {
		for _, v := range errs {
			slog.Error("app fatal error", "comment", comment, "error", v)
		}
		os.Exit(int(ExitError))
	}
}
