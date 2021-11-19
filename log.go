package main

// Simple logging functions *with color*

import (
	"fmt"

	"github.com/fatih/color"
)

func warn(text string) {
	fmt.Printf("%v%v%v%s\n", color.HiBlackString("("), color.YellowString("-"), color.HiBlackString(") "), text)
}

func info(text string) {
	fmt.Printf("%v%v%v%s\n", color.HiBlackString("("), color.GreenString("+"), color.HiBlackString(") "), text)
}

func logerror(text string) {
	fmt.Printf("%v%v%v%s\n", color.HiBlackString("("), color.RedString("!"), color.HiBlackString(") "), text)
}
