package main

// import (
// 	"testing"
// 	"io"
// 	"os"
// )

// func TestGetLogFile(t *testing.T) {
// 	type Test struct {
// 		path string
// 		stderr io.Writer
// 	}
// 	var tests = []Test {
// 		{"", os.Stderr}, 
// 		{" ", os.Stderr}, 
// 		{"/var/tmp/atat.log", os.Stdout},
// 	}

// 	for i, test := range tests {
// 		LogFile := getLogFile(test.path)
// 		if LogFile != test.stderr {
// 			t.Errorf("test %d: %s", i, test.path)
// 		}
// 	}
// }
