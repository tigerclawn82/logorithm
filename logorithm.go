package logger

import (
	"fmt"
	"io"
	"log"
	"time"
)

type Logger struct {
	*log.Logger
	sequenceNumber int
	verbose        bool
	formatString   string
}

const BASE_FORMAT_STR = "%%s[origin software=\"%s\" swVersion=\"%s\" x-program=\"%s\" x-sequence=\"%%d\" x-pid=\"%d\" x-severity=\"%%s\" x-timestamp=\"%%s\"] %%s"

func New(stream io.Writer, verbose bool, software string, version string, program string, pid int) *Logger {
	return &Logger{
		Logger:         log.New(stream, "", 0),
		sequenceNumber: 0,
		verbose:        verbose,
		formatString:   fmt.Sprintf(BASE_FORMAT_STR, software, version, program, pid),
	}
}

func (l *Logger) Log(severity string, msg string, vargs ...interface{}) {
	l.Printf(l.formatString,
		severity,
		l.sequenceNumber,
		severity,
		time.Now().Format(time.RFC3339Nano),
		fmt.Sprintf(msg, vargs...),
	)
	l.sequenceNumber++
}

func (l *Logger) Emerg(msg string, vargs ...interface{}) {
	l.Log("EMERG", msg, vargs...)
}
func (l *Logger) Alert(msg string, vargs ...interface{}) {
	l.Log("ALERT", msg, vargs...)
}
func (l *Logger) Critical(msg string, vargs ...interface{}) {
	l.Log("CRIT", msg, vargs...)
}
func (l *Logger) Error(msg string, vargs ...interface{}) {
	l.Log("ERR", msg, vargs...)
}
func (l *Logger) Warning(msg string, vargs ...interface{}) {
	l.Log("WARN", msg, vargs...)
}
func (l *Logger) Notice(msg string, vargs ...interface{}) {
	l.Log("NOTICE", msg, vargs...)
}
func (l *Logger) Info(msg string, vargs ...interface{}) {
	l.Log("INFO", msg, vargs...)
}
func (l *Logger) Debug(msg string, vargs ...interface{}) {
	if l.verbose {
		l.Log("DEBUG", msg, vargs...)
	}
}
