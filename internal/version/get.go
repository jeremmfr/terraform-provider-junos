package version

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var current string

// Get the version.
func Get() string {
	return strings.TrimSpace(current)
}
