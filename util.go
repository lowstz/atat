package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"unicode"
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

// Verify the validity of ISBN number
func  isValidIsbn13(isbn string) bool {
	var (
		i int
		check int = 0
	)
	if len(isbn) != 13 {
		return false
	}
	for _,r := range isbn {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	isbnArray := strings.Split(isbn, "")
	for i=0; i<13; i+=2 {
		digits, _ := strconv.Atoi(isbnArray[i])
		check += digits
	}
	for i=1; i<12; i+=2 {
		digits, _ := strconv.Atoi(isbnArray[i])
		check += 3*digits
	}
	return check%10 == 0
}

// Check error and panic it.
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
