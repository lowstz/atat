package main

import (
	"testing"
	"path/filepath"
	"os"
	"os/user"
)


func TestExpandPath(t *testing.T) {
	paths := []string{
		"foo",
		"~",
		"~foo",
		"~/foo",
		"foo/~/",
		"$HOME/foo",
	}

	cwd := os.Getenv("PWD")
	u, _:= user.Current()
	home := u.HomeDir
	expanded := []string{
		filepath.Join(cwd, "foo"),  // ./foo
		home,                       // /home/current/
		filepath.Join(cwd, "~foo"),   // ./~foo
		filepath.Join(home, "foo"),   // /home/current/foo
		filepath.Join(cwd, "foo/~/"), // ./foo/~/
		filepath.Join(home, "foo"),   // /home/current/foo
		filepath.Join(home, "foo"),   // /home/current/foo
	}

	for i, path := range paths {
		res := ExpandPath(path)
		if res != expanded[i] {
			t.Errorf("%d. Expected '%s' => '%s', got '%s'", i, paths[i], expanded[i], res)
		}
	}
}

func TestIsFileExists(t *testing.T) {
	type Test struct {
		path string
		ispath bool
	}

	var tests = []Test{
		{".", false},
		{"..", false},
		{"./conf/config.conf", true},
		{"/etc/passwd", true},
		{"~/.bashrc", true},
		{"~/.bashr~c", false},
	}

	for i, test := range tests {
		exists, _ := IsFileExists(test.path)
		if exists != test.ispath {
			t.Errorf("test %d: %s %#v != %#v", i, test.path , exists, test.ispath)
		}
	}
}


func TestIsValidIsbn13(t *testing.T) {
	type Test struct {
		isbn string
		isValid bool
	}

	var tests = []Test {
		{"9787302274759", true},
		{"9787515301723", true},
		{"9787301199800", true},
		{"9787302273158", true},
		{"af3jf", false},
		{"bbbbbbbbbbbbb", false},
		{"978730227315X", false},
		{"9787238082399", false},
	}

	for i, test := range tests {
		isvalid := isValidIsbn13(test.isbn)
		if isvalid != test.isValid {
			t.Errorf("test %d: %s %#v != %#v", i, test.isbn , isvalid, test.isValid)
		}
	}
}
