package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Fatal("failed to read config.json file:", err)
	}

	var config internal.Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("failed to parse config.json file:", err)
	}

	count, err := config.SetSymLink()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Printf("%s[TOTAL SYMLINK]%s: %d\n", internal.CYAN_COLOR, internal.RESET_COLOR, count)
}
