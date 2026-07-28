package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	CURRENT_VERSION uint8 = 1
	GREEN_COLOR           = "\033[32m"
	RED_COLOR             = "\033[31m"
	RESET_COLOR           = "\033[0m"
)

type Path struct {
	Name     string `json:"name"`
	Original string `json:"original_path"`
	Target   string `json:"target_path"`
	Desc     string `json:"description"`
}

func (p *Path) Sanitize() {
	// original
	if p.Original == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(69)
		}
		p.Original = home
	} else if strings.HasPrefix(p.Original, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(69)
		}

		p.Original = filepath.Join(home, p.Original[2:])
	}

	// target
	if p.Target == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(69)
		}
		p.Target = home
	} else if strings.HasPrefix(p.Target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(69)
		}

		p.Target = filepath.Join(home, p.Target[2:])
	}
}

type Config struct {
	SymPaths []Path `json:"sym_paths"`
	Author   string `json:"author"`
	Version  uint8  `json:"version"`
}

func (c Config) IsExist(path string) error {
	stat, exist := os.Lstat(path)
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
		p.Sanitize()
		basePath := filepath.Base(p.Original)
		fullPath := filepath.Join(p.Target, basePath)
		err := c.IsExist(fullPath)
		if err != nil {
			return fmt.Errorf(
				"%s[ERROR]%s: %s : %w",
				RED_COLOR,
				RESET_COLOR,
				p.Name,
				err,
			)
		}
		err = os.Symlink(p.Original, fullPath)
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

func main() {
	args := os.Args

	if len(args) > 2 || len(args) == 1 {
		fmt.Fprintf(
			os.Stderr,
			"%s[ERROR]%s: wrong arguments passed, correct arguments: <program> <config-json-file>\n",
			RED_COLOR,
			RESET_COLOR,
		)
		return
	}

	jsonFile := args[1]
	content, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal("failed to read config.json file:", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("failed to parse config.json file:", err)
	}

	err = config.SetSymLink()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("TOTAL SYMLINK: %d\n", len(config.SymPaths))
}
