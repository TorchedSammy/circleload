package log

// Simple logging functions *with color*

import (
	"fmt"

	"github.com/fatih/color"
)

func Warn(text ...interface{}) {
	fmt.Printf("%v%v%v%s\n", color.HiBlackString("("), color.YellowString("-"), color.HiBlackString(") "), fmt.Sprintln(text...))
}

func Info(text ...interface{}) {
	fmt.Printf("%v%v%v%s\n", color.HiBlackString("("), color.GreenString("+"), color.HiBlackString(") "), fmt.Sprintln(text...))
}

func Error(text ...interface{}) {
	fmt.Printf("%v%v%v%s\n", color.HiBlackString("("), color.RedString("!"), color.HiBlackString(") "), fmt.Sprintln(text...))
}
