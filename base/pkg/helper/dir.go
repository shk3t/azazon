package helper

import (
	"os"
	"path/filepath"
)

func GetwdCdBack(unwantedDirChain ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	chainLen := len(unwantedDirChain)
	for i := range chainLen {
		_, wdTail := filepath.Split(wd)
		if wdTail == unwantedDirChain[chainLen-i-1] {
			wd = filepath.Dir(wd)
		}
	}

	return wd, nil
}