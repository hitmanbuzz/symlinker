package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Sanitizer struct {
	path *Path
}

func NewSanitizer(path *Path) *Sanitizer {
	return &Sanitizer{path}
}

// handle home directory
//
// eg: ~/dotfiles/config.toml -> /home/username/dotfiles/config.toml
func (s *Sanitizer) sanitize_home() error {
	// original
	if s.path.Original == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		s.path.Original = home
	} else if strings.HasPrefix(s.path.Original, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		s.path.Original = filepath.Join(home, s.path.Original[2:])
	}

	// target
	if s.path.Target == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		s.path.Target = home
	} else if strings.HasPrefix(s.path.Target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		s.path.Target = filepath.Join(home, s.path.Target[2:])
	}
	return nil
}

// handle current directory
//
// eg: ./config.toml
func (s *Sanitizer) sanitize_current() error {
	if s.path.Original == "." || strings.HasPrefix(s.path.Original, "./") {
		// shouldn't have any error
		lastSlash := strings.LastIndex(ConfigPath, "/")
		if lastSlash == -1 {
			return fmt.Errorf("`/` not found in the path")
		}
		s.path.Original = ConfigPath[:lastSlash]
	}

	// target
	if s.path.Target == "." || strings.HasPrefix(s.path.Target, "./") {
		// shouldn't have any error
		lastSlash := strings.LastIndex(ConfigPath, "/")
		if lastSlash == -1 {
			return fmt.Errorf("`/` not found in the path")
		}
		s.path.Target = ConfigPath[:lastSlash]
	}

	return nil
}

// handle previous directory
//
// eg: ../../config.toml
func (s *Sanitizer) sanitize_previous() error {
	return nil
}

func (s *Sanitizer) Sanitize() error {
	s.path.Entity = filepath.Base(s.path.Original)

	err := s.sanitize_home()
	if err != nil {
		return err
	}

	err = s.sanitize_current()
	if err != nil {
		return err
	}

	err = s.path.isExist()
	if err != nil {
		return fmt.Errorf(
			"%s[ERROR]%s: %s : %w",
			RED_COLOR,
			RESET_COLOR,
			s.path.Name,
			err,
		)
	}

	return nil
}
