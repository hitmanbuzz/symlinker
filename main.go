package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var ConfigJsonPath string

const (
	CURRENT_VERSION uint8 = 1
	GREEN_COLOR           = "\033[32m"
	YELLOW_COLOR          = "\033[33m"
	RED_COLOR             = "\033[31m"
	CYAN_COLOR            = "\033[36m"
	BLUE_COLOR            = "\033[34m"
	RESET_COLOR           = "\033[0m"
)

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

	ConfigJsonPath = args[1]
	content, err := os.ReadFile(ConfigJsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read `%s` file: %v\n", ConfigJsonPath, err)
		return
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse `%s` file: %v\n", ConfigJsonPath, err)
		return
	}

	s := NewSymlink(&config)
	count, err := s.SetSymLink()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s[TOTAL SYMLINK]%s: %d\n", CYAN_COLOR, RESET_COLOR, count)
}
