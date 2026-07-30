package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

type Path struct {
	Name     string `json:"name"`
	Original string `json:"original_path"`
	Target   string `json:"target_path"`
	Desc     string `json:"description"`
	Entity   string
}

type Config struct {
	SymPaths []Path `json:"sym_paths"`
	Author   string `json:"author"`
	Version  uint8  `json:"version"`
}

func (p Path) isExist() (error, bool) {
	err := IsExist(p.Original)
	if err != nil {
		return fmt.Errorf("original_path: `%s` doesn't exist", p.Original), false
	}

	targetFullPath := filepath.Join(p.Target, p.Entity)
	stat, exist := os.Lstat(targetFullPath)
	if exist == nil {
		if stat.Mode()&os.ModeSymlink != 0 {
			err, isSame := IsSameSymlink(p.Original, targetFullPath, p.Entity)
			return err, isSame
		} else {
			return fmt.Errorf("a non-symlink path already exist"), false
		}
	}
	return nil, false
}

func (c Config) SetSymLink() (uint, error) {
	if c.Version != CURRENT_VERSION {
		return 0, fmt.Errorf(
			"%s[ERROR]%s: config version is %d but the program support version %d only",
			RED_COLOR,
			RESET_COLOR,
			c.Version,
			CURRENT_VERSION,
		)
	}

	var symlinkCount uint = 0

	for _, p := range c.SymPaths {
		s := NewSanitizer(&p)
		err := s.Sanitize()
		if err != nil {
			return symlinkCount, err
		}

		err, exist := p.isExist()
		if err != nil {
			return symlinkCount, fmt.Errorf(
				"%s[ERROR]%s: %s : %w",
				RED_COLOR,
				RESET_COLOR,
				s.path.Name,
				err,
			)
		}

		if exist {
			continue
		}

		err = os.Symlink(p.Original, filepath.Join(p.Target, p.Entity))
		if err != nil {
			return symlinkCount, fmt.Errorf(
				"%s[ERROR]%s: %w",
				RED_COLOR,
				RESET_COLOR,
				err,
			)
		}
		symlinkCount++
		fmt.Printf("%s[SUCCESS]%s: %s\n", GREEN_COLOR, RESET_COLOR, p.Name)
	}
	return symlinkCount, nil
}
