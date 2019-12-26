package base

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

var level = INFO

const (
	_ = iota
	DEBUG
	INFO
)

func fmtLogLevel(logLevel int) string {
	if logLevel < INFO {
		return "DEBUG"
	}
	return "INFO"
}

func Debug(format string, a ...interface{}) {
	logFull(DEBUG, "", "", format, a...)
}

func Info(format string, a ...interface{}) {
	logFull(INFO, "", "", format, a...)
}

func DebugFull(repoID string, runID string, format string, a ...interface{}) {
	logFull(DEBUG, repoID, runID, format, a...)
}

func SetLogFlags(flags int) {
	log.SetFlags(flags)
}

func SetLogLevel(newLevel int) {
	level = newLevel
	log.SetFlags(0)
}

func SetLogFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot open log file %v", path))
	}
	log.SetOutput(f)
	return nil
}

func logFull(logLevel int, runID string, repoID string, format string, a ...interface{}) {
	// fields are
	// time, log_level, run_id, repo_id, message
	if runID == "" {
		runID = "Null"
	}
	if repoID == "" {
		repoID = "Null"
	}

	if logLevel >= level {
		line := fmt.Sprintf("%v %s %s %s %s",
			time.Now().Format(time.RFC3339),
			fmtLogLevel(logLevel),
			runID,
			repoID,
			fmt.Sprintf(format, a...))
		log.Println(string(line))
	}
}
