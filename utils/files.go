package utils

import (
	"errors"
	"os"
)

func CreateDirectories(dirs []string) error {
	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}
