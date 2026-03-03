package fsys

import (
	"errors"
	"io"
	"os"
)

func PreviewFile(s string, l int) ([]byte, error) {
	f, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, l)
	n, err := f.Read(buf)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}

		return nil, err
	}

	return buf[:n], nil
}

func DirEmpty(s string) (bool, error) {
	f, err := os.Open(s)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

func ExistsStat(s string) (bool, os.FileInfo, error) {
	i, err := os.Stat(s)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil, nil
		}

		return false, nil, err
	}

	return true, i, nil
}
