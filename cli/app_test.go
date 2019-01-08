package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewMyApp(t *testing.T) {
	app := NewMyApp()
	if app.stream != nil {
		t.Fatalf("stream should be initialized with nil")
	}
}

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app.stream != nil {
		t.Fatalf("stream should be initialized with nil")
	}
}

func TestAction(t *testing.T) {
	app := NewApp()
	tests := []struct {
		path string
		size int64
	}{
		{"fileA", 10 * 1024 * 1024},
		{"fileB", 10*1024 + 100},
		{"fileC", 512},
	}

	tmpfiles := []string{}
	for _, test := range tests {
		path := fmt.Sprintf("%s/%s", os.TempDir(), test.path)
		createFile(path, test.size)
		tmpfiles = append(tmpfiles, path)
		defer func() { os.Remove(path) }()
	}
	err := app.Run(append([]string{"-q"}, tmpfiles...))
	if err != nil {
		t.Fatalf("Error happened: %s", err.Error())
	}
	for _, f := range tmpfiles {
		_, perr := os.Stat(f)
		if perr != nil {
			if _, ok := perr.(*os.PathError); !ok {
				t.Fatalf("Filetest for %s failed with non-PathError: %s", f, perr.Error())
			}
		} else {
			t.Fatalf("File %s is not removed", f)
		}
	}
}

func TestAction_PathError(t *testing.T) {
	app := NewApp()
	tests := []struct {
		path string
		size int64
	}{
		{"fileD", 512},
	}

	tmpfiles := []string{}
	for _, test := range tests {
		path := fmt.Sprintf("%s/%s", os.TempDir(), test.path)
		// not creating file
		tmpfiles = append(tmpfiles, path)
		defer func() { os.Remove(path) }()
	}
	err := app.Run(append([]string{"-q"}, tmpfiles...))
	if err == nil {
		t.Fatalf("Error did not happen (it should)")
	}
}

func createFile(path string, size int64) {
	file, err := os.Create(path)
	if err != nil {
		panic("Failed to create file " + path + " " + err.Error())
	}
	_, err = file.WriteAt([]byte("a"), size-1)
	if err != nil {
		panic("Failed to write file " + path + " " + err.Error())
	}
}
