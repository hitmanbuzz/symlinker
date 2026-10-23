package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SymPaths []*Path `json:"sym_paths"`
	Version  float32 `json:"version"`
}

type Path struct {
	Name     string `json:"name"`
	Original string `json:"original_path"`
	Target   string `json:"target_path"`
	Base     string
}

func (p *Path) IsExist() (bool, error) {
	originalInfo, err := os.Stat(p.Original)
	if err != nil {
		return false, fmt.Errorf(
			"original path doesn't exist: `%s`",
			p.Original,
		)
	}

	fullTargetPath := filepath.Join(p.Target, p.Base)

	targetLinfo, err := os.Lstat(fullTargetPath)
	if err != nil {
		return false, nil
	}

	if p.IsSymlink(targetLinfo) {
		targetSinfo, err := os.Stat(fullTargetPath)
		if err != nil {
			if os.IsNotExist(err) {
				return false, fmt.Errorf(
					"target path already exist with a broken symlink: `%s`",
					fullTargetPath,
				)
			} else {
				return false, err
			}
		}

		if os.SameFile(originalInfo, targetSinfo) {
			return true, nil
		} else {
			oPath, _ := os.Readlink(fullTargetPath)
			return false, fmt.Errorf(
				"target already exist but it symlink to different original path: `%s`",
				oPath,
			)
		}
	} else {
		if targetLinfo.IsDir() {
			return true, fmt.Errorf(
				"target already exist but as a directory instead of symlink: `%s`",
				p.Target,
			)
		} else if !targetLinfo.IsDir() {
			return true, fmt.Errorf(
				"target already exist but as a file instead of symlink: `%s`",
				p.Target,
			)
		}
	}

	return false, nil
}

func (p *Path) Sanitize() error {
	p.Base = filepath.Base(p.Original)
	sanitize(&p.Original)
	sanitize(&p.Target)
	return nil
}

func (p Path) IsSymlink(info os.FileInfo) bool {
	if info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}

func sanitize(p *string) error {
	if *p == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		*p = home
	} else if strings.HasPrefix(*p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		*p = filepath.Join(home, (*p)[2:])
	} else {
		absPath, err := filepath.Abs(*p)
		if err != nil {
			return err
		}
		*p = absPath
	}

	return nil
}
