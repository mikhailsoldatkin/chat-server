// Package my_logger временное решение для цветных логов
package my_logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Fatal prints a formatted message in red color to the standard output and exits the program.
func Fatal(format string, args ...interface{}) {
	red := color.New(color.FgRed).PrintfFunc()
	red(format, args...)
	fmt.Println()
	os.Exit(1)
}

// Info prints a formatted message in green color to the standard output.
func Info(format string, args ...interface{}) {
	green := color.New(color.FgGreen).PrintfFunc()
	green(format, args...)
	fmt.Println()
}
