package file

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// WithUserHomeDir is a helper that takes in a slice of paths and concats it
// with the home directory
func WithUserHomeDir(paths ...string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	fullPath := []string{home}

	fullPath = append(fullPath, paths...)

	return filepath.Join(fullPath...), nil
}
