package cmd

import "fmt"

//goland:noinspection SpellCheckingInspection
func vprintf(level int, format string, args ...interface{}) {
	if verbose >= level {
		fmt.Printf(format, args...)
	}
}
