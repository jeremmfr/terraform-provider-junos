package utils

import (
	"fmt"
	"os"

	balt "github.com/jeremmfr/go-utils/basicalter"
)

func ReplaceTildeToHomeDir(path *string) error {
	if balt.CutPrefixInString(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to read user home directory: %w", err)
		}
		*path = homeDir + *path
	}

	return nil
}
