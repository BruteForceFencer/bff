// Package version provides access to the current version of BFF.
package version

import (
	"fmt"
)

// Version is the current version of BFF.
const Version = "0.4.0"

// PrintVersion prints the version of BFF to standard output.
func PrintVersion() {
	fmt.Println("BFF core version", Version)
	fmt.Println("Copyright (C) James Hall 2015.")
}
