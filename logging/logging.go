package logging

import (
	"log"
	"io/ioutil"
	_ "os"
)

// Discard is a Log that just discards its messages.
var Discard = log.New(ioutil.Discard, "", log.LstdFlags)

type Log struct {
	logger *log.Logger
}

func (lg *Log)SetLogger(logger *log.Logger) {
	lg.logger = logger
}

func (lg *Log)Logger() *log.Logger {
	if lg.logger == nil {
				lg.logger = log.New(ioutil.Discard, "", log.LstdFlags)
		//lg.logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return lg.logger
}

