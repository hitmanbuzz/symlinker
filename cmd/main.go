package main

import (
	"encoding/json"
	"fmt"
	"os"
	"symlinker/internal"
)

func main() {
	args := os.Args

	if len(args) > 2 || len(args) == 1 {
		fmt.Fprintf(
			os.Stderr,
			"%s[ERROR]%s: wrong arguments passed, correct arguments: <program> <config-json-file>\n",
			internal.RED_COLOR,
			internal.RESET_COLOR,
		)
		return
	}

	internal.ConfigPath = args[1]
	content, err := os.ReadFile(internal.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "faild to read `%s` file: %v\n", internal.ConfigPath, err)
		return
	}

	var config internal.Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse `%s` file: %v\n", internal.ConfigPath, err)
		return
	}

	count, err := config.SetSymLink()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Printf("%s[TOTAL SYMLINK]%s: %d\n", internal.CYAN_COLOR, internal.RESET_COLOR, count)
}
