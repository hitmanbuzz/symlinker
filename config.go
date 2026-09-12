package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	SymPaths []*Path `json:"sym_paths"`
	Version  uint8   `json:"version"`
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
		if targetLinfo.IsDir() && originalInfo.IsDir() {
			return true, fmt.Errorf(
				"target already exist as a directory: `%s`",
				p.Target,
			)
		} else if !targetLinfo.IsDir() && !originalInfo.IsDir() {
			return true, fmt.Errorf(
				"target already exist as a file: `%s`",
				p.Target,
			)
		}
	}

	return false, nil
}

func (p *Path) Sanitize() error {
	p.Base = filepath.Base(p.Original)

	// original
	if p.Original == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		p.Original = home
	} else if strings.HasPrefix(p.Original, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		p.Original = filepath.Join(home, p.Original[2:])
	} else {
		absPath, err := filepath.Abs(p.Original)
		if err != nil {
			return err
		}
		p.Original = absPath
	}

	// target
	if p.Target == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		p.Target = home
	} else if strings.HasPrefix(p.Target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		p.Target = filepath.Join(home, p.Target[2:])
	} else {
		absPath, err := filepath.Abs(p.Target)
		if err != nil {
			return err
		}
		p.Target = absPath
	}

	return nil
}

func (p Path) IsSymlink(info os.FileInfo) bool {
	if info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}
