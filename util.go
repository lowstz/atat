package main

import (
	"errors"
	"os"
	"os/user"
	"strconv"
	"strings"
	"unicode"
	"path/filepath"
)

// expand file path to absolute path
func ExpandPath(path string) string {
	if path[0] == '~' {
		sep := strings.Index(path, string(os.PathSeparator))
		if sep < 0 {
			sep = len(path)
		}
		var err error
		var u *user.User
		username := path[1:sep]
		if len(username) == 0 {
			u, err = user.Current()
		} else {
			u, err = user.Lookup(username)
		}
		if err == nil {
			path = filepath.Join(u.HomeDir, path[sep:])
		}
	}
	path = os.ExpandEnv(path)
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}


// Check file is Exists.
func IsFileExists(path string) (bool, error) {
	expandPath  := ExpandPath(path)
	stat, err := os.Stat(expandPath)
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
