package logger

import (
	"os"

	l "github.com/charmbracelet/log"
)

var DefaultLogger = l.NewWithOptions(os.Stderr, l.Options{
	Level:           l.InfoLevel,
	ReportTimestamp: false,
})

var VerboseLogger = l.NewWithOptions(os.Stderr, l.Options{
	Level:           l.DebugLevel,
	ReportTimestamp: false,
})
