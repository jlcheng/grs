package base

import "fmt"

var level = INFO

const (
	DEBUG = iota
	INFO  = iota
)

func Debug(format string, a ...interface{}) {
	if level <= DEBUG {
		fmt.Printf("[DEBUG] %v\n", fmt.Sprintf(format, a...))
	}
}

func Info(format string, a ...interface{}) {
	if level <= INFO {
		fmt.Printf("[INFO] %v\n", fmt.Sprintf(format, a...))
	}
}

func SetLogLevel(newLevel int) {
	level = newLevel
}
