package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type Symlink struct {
	config *Config
}

func NewSymlink(config *Config) *Symlink {
	return &Symlink{
		config: config,
	}
}

func (s *Symlink) SetSymLink() (int, error) {
	if s.config.Version != CURRENT_VERSION {
		return 0, fmt.Errorf(
			"%s[ERROR]%s: config version is %d but the program support version %d only",
			RED_COLOR,
			RESET_COLOR,
			s.config.Version,
			CURRENT_VERSION,
		)
	}

	count := 0

	for _, p := range s.config.SymPaths {
		err := p.Sanitize()
		if err != nil {
			return count, err
		}

		exist, err := p.IsExist()
		if err != nil {
			return count, fmt.Errorf(
				"%s[ERROR]%s: %s : %w",
				RED_COLOR,
				RESET_COLOR,
				p.Name,
				err,
			)
		}

		if exist && err == nil {
			fmt.Printf("%s[DUPLICATE]%s: %s\n", YELLOW_COLOR, RESET_COLOR, p.Name)
			continue
		}

		err = os.Symlink(p.Original, filepath.Join(p.Target, p.Base))
		if err != nil {
			return count, fmt.Errorf(
				"%s[ERROR]%s: %w",
				RED_COLOR,
				RESET_COLOR,
				err,
			)
		}
		count++
		fmt.Printf("%s[SUCCESS]%s: %s\n", GREEN_COLOR, RESET_COLOR, p.Name)
	}

	return count, nil
}
