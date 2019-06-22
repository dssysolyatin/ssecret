package filepathu

import (
	"os/user"
	"path/filepath"
)

func Abs(path string) (string, error) {
	if path[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return "", err
		}

		path = filepath.Join(u.HomeDir, path[1:])
	}

	return filepath.Clean(path), nil
}
