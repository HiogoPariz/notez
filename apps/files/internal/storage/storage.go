package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/HiogoPariz/notez/apps/files/internal/config"
)

// ErrInvalidPath is returned when the requested key is not a safe file name.
var ErrInvalidPath = errors.New("invalid path")

// validPath only allows simple keys (hex hashes, slugs) so the URL parameter
// can never escape the storage directory.
var validPath = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func dir() string {
	return config.Get("FILES_DIR", ".files")
}

// resolve validates the key and returns the full path of its .json file.
func resolve(path string) (string, error) {
	if !validPath.MatchString(path) {
		return "", fmt.Errorf("%w: %q", ErrInvalidPath, path)
	}
	return filepath.Join(dir(), path+".json"), nil
}

func ReadFile(path string) (string, error) {
	full, err := resolve(path)
	if err != nil {
		return "", err
	}

	file, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func WriteFile(path string, content string) error {
	full, err := resolve(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir(), 0o755); err != nil {
		return err
	}

	return os.WriteFile(full, []byte(content), 0o644)
}

func DeleteFile(path string) error {
	full, err := resolve(path)
	if err != nil {
		return err
	}

	return os.Remove(full)
}
