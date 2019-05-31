package base

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
)

var level = INFO

const (
	DEBUG = iota
	INFO  = iota
)

func Debug(format string, a ...interface{}) {
	if level <= DEBUG {
		log.Printf("[DEBUG] %v\n", fmt.Sprintf(format, a...))
	}
}

func Info(format string, a ...interface{}) {
	if level <= INFO {
		log.Printf("[INFO] %v\n", fmt.Sprintf(format, a...))
	}
}

func SetLogLevel(newLevel int) {
	level = newLevel
}

func SetLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot open log file %v", path))
	}
	log.SetOutput(f)
	return nil
}
