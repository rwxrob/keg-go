package keg

import (
	"os"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func notExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}
