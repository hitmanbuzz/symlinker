package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/sumdb/dirhash"
)

var ConfigPath string

const (
	CURRENT_VERSION uint8 = 1
	GREEN_COLOR           = "\033[32m"
	YELLOW_COLOR          = "\033[33m"
	RED_COLOR             = "\033[31m"
	CYAN_COLOR            = "\033[36m"
	RESET_COLOR           = "\033[0m"
)

func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func IsSameSymlink(original string, target string, entity string) (error, bool) {
	oriStat, err := os.Stat(original)
	if err != nil {
		return err, false
	}

	targetInfo, err := os.Stat(target)
	if err != nil {
		return err, false
	}

	if targetInfo.IsDir() && oriStat.IsDir() {
		oriHash, err := dirhash.HashDir(original, "", dirhash.DefaultHash)
		if err != nil {
			return err, false
		}

		realTargetPath, err := filepath.EvalSymlinks(target)
		if err != nil {
			return err, false
		}

		targetHash, err := dirhash.HashDir(realTargetPath, "", dirhash.DefaultHash)
		if err != nil {
			return err, false
		}

		if oriHash == targetHash {
			fmt.Printf("%s[DUPLICATE]%s: %s\n", YELLOW_COLOR, RESET_COLOR, entity)
			return nil, true
		} else {
			return fmt.Errorf("a different symlink with same name already exist"), false
		}
	} else if !targetInfo.IsDir() && !oriStat.IsDir() {
		oriHash, err := HashFile(original)
		if err != nil {
			return err, false
		}

		realTargetPath, err := filepath.EvalSymlinks(target)
		if err != nil {
			return err, false
		}

		targetHash, err := HashFile(realTargetPath)
		if err != nil {
			return err, false
		}

		if oriHash == targetHash {
			fmt.Printf("%s[DUPLICATE]%s: %s\n", YELLOW_COLOR, RESET_COLOR, entity)
			return nil, true
		} else {
			return fmt.Errorf("a different symlink with same name already exist"), false
		}
	} else {
		return fmt.Errorf("a symlink of different format already exist"), false
	}
}
