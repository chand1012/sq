package logger

import (
	"os"

	l "github.com/charmbracelet/log"
)

func HandlePanic(log *l.Logger, err error, verbose bool) {
	log.Error(err.Error())
	if verbose {
		panic(err)
	}
	os.Exit(1)
}
