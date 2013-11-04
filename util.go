package main

import (
	"os"
	"errors"
	"log"
)


// Check file is Exists.
func isFileExists(path string) (bool, error) {
        stat, err := os.Stat(path)
        if err == nil {
                if stat.Mode()&os.ModeType == 0 {
                        return true, nil
                }
                return false, errors.New(path + " exists but is not regular file")
        }
        if os.IsNotExist(err) {
                return false, nil
        }
        return false, err
}

// Check error
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}











