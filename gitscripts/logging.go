package gitscripts

import "fmt"

var level int

const (
	DEBUG = iota
	INFO = iota
)

func Debug(format string, a ...interface{}) {
	if (level <= DEBUG) {
		fmt.Printf(format, a)
		fmt.Println()
	}
}

func Info(format string, a ...interface{}) {
	if (level <= INFO) {
		fmt.Printf(format, a)
		fmt.Println()
	}
}

func SetLogLevel(newLevel int) {
	level = newLevel
}