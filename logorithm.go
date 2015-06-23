package logorithm

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"
)

type Logger interface {
	Emerg(string, ...interface{})
	Alert(string, ...interface{})
	Critical(string, ...interface{})
	Error(string, ...interface{})
	Warning(string, ...interface{})
	Notice(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
}

type L struct {
	*log.Logger
	sync.Mutex
	sequenceNumber int
	verbose        bool
	formatString   string
}

const BASE_FORMAT_STR = "%%s[origin software=\"%s\" swVersion=\"%s\" x-program=\"%s\" x-sequence=\"%%d\" x-pid=\"%d\" x-severity=\"%%s\" x-timestamp=\"%%s\"] %%s"

func New(stream io.Writer, verbose bool, software string, version string, program string, pid int) *L {
	return &L{
		Logger:         log.New(stream, "", 0),
		sequenceNumber: 0,
		verbose:        verbose,
		formatString:   fmt.Sprintf(BASE_FORMAT_STR, software, version, program, pid),
	}
}

func Instance(stream io.Writer, verbose bool, software string, version string, program string, pid int, logFmt string) *L {
	if logFmt != "" {
		logFmt = fmt.Sprintf(logFmt, software, version, program, pid)
	} else {
		logFmt = fmt.Sprintf(BASE_FORMAT_STR, software, version, program, pid)
	}
	return &L{
		Logger:         log.New(stream, "", 0),
		sequenceNumber: 0,
		verbose:        verbose,
		formatString:   logFmt,
	}
}

func (l *L) increment() int {
	l.Lock()
	defer l.Unlock()

	l.sequenceNumber++
	return l.sequenceNumber
}

func (l *L) log(severity string, msg string, vargs ...interface{}) {
	l.Printf(l.formatString,
		severity,
		l.increment(),
		severity,
		time.Now().Format(time.RFC3339Nano),
		fmt.Sprintf(msg, vargs...),
	)
}

func (l *L) Emerg(msg string, vargs ...interface{}) {
	l.log("EMERG", msg, vargs...)
}
func (l *L) Alert(msg string, vargs ...interface{}) {
	l.log("ALERT", msg, vargs...)
}
func (l *L) Critical(msg string, vargs ...interface{}) {
	l.log("CRIT", msg, vargs...)
}
func (l *L) Error(msg string, vargs ...interface{}) {
	l.log("ERR", msg, vargs...)
}
func (l *L) Warning(msg string, vargs ...interface{}) {
	l.log("WARN", msg, vargs...)
}
func (l *L) Notice(msg string, vargs ...interface{}) {
	l.log("NOTICE", msg, vargs...)
}
func (l *L) Info(msg string, vargs ...interface{}) {
	l.log("INFO", msg, vargs...)
}
func (l *L) Debug(msg string, vargs ...interface{}) {
	if l.verbose {
		l.log("DEBUG", msg, vargs...)
	}
}
