package internal

import (
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

func (s *Sanitizer) sanitize_path() error {
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
	} else {
		s.path.Original = filepath.Dir(ConfigPath)
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
	} else {
		s.path.Target = filepath.Dir(ConfigPath)
	}

	return nil
}

func (s *Sanitizer) Sanitize() error {
	s.path.Entity = filepath.Base(s.path.Original)

	err := s.sanitize_path()
	if err != nil {
		return err
	}

	return nil
}
