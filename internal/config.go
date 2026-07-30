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

// FIX: compare the data hash of the path with the original path data hash
//
// This make sure that the current existing symlink is the same as the original or not
func (p Path) isExist() error {
	fullPath := filepath.Join(p.Target, p.Entity)
	stat, exist := os.Lstat(fullPath)
	if exist == nil {
		if stat.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("a symlink path already exist")
		} else {
			return fmt.Errorf("a non-symlink path already exist")
		}
	}
	return nil
}

func (c Config) SetSymLink() error {
	if c.Version != CURRENT_VERSION {
		return fmt.Errorf(
			"%s[ERROR]%s: config version is %d but the program support version %d only",
			RED_COLOR,
			RESET_COLOR,
			c.Version,
			CURRENT_VERSION,
		)
	}

	for _, p := range c.SymPaths {
		s := NewSanitizer(&p)
		err := s.Sanitize()
		if err != nil {
			return err
		}
		err = os.Symlink(p.Original, p.Entity)
		if err != nil {
			return fmt.Errorf(
				"%s[ERROR]%s: %w",
				RED_COLOR,
				RESET_COLOR,
				err,
			)
		}
		fmt.Printf("%s[SUCCESS]%s: %s\n", GREEN_COLOR, RESET_COLOR, p.Name)
	}
	return nil
}
