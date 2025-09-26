package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type Local struct {
	maxFileSize int
	basePath    string
}

func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{basePath: p, maxFileSize: maxSize}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	// / get the full path for the file
	fp := l.fullPath(path)
	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}
	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil { // means fp file path has no stats to give
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Unable to get file info: %w", err)
	}
	// Create a new file at path
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()
	// Write the contents
	// Copy copies from src to dst until either EOF is
	// reached on src or an error occurs.
	// It returns the number of bytes copied and the first error encountered
	// while copying, if any.
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to the filepath: %w", err)
	}
	return nil
}

func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}
