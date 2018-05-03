package grsio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copy(src string, dest string) error {
	srcstat, e := os.Stat(src)
	if srcstat == nil {
		return e
	}
	if srcstat.IsDir() {
		return copydir(src, dest)
	}
	return copyfile(src, dest)
}

func copydir(src string, dest string) error {
	// src must exist
	srcdir, e := os.Open(src)
	if srcdir == nil {
		return e
	}
	defer srcdir.Close()

	// dest MUST NOT exist
	_, e = os.Stat(dest)
	if !os.IsNotExist(e) {
		return errors.New(fmt.Sprintf("destination %v already exists", dest))
	}

	// create dest
	e = os.Mkdir(dest, 0777)
	if e != nil {
		return e
	}

	// identify child nodes from src
	files, e := srcdir.Readdir(-1)
	if e != nil {
		return e
	}
	e = srcdir.Close()
	if e != nil {
		return e
	}

	// copy child nodes
	for _, file := range files {
		schild := filepath.Join(src, file.Name())
		dchild := filepath.Join(dest, file.Name())
		copy(schild, dchild)
	}
	return nil
}

func copyfile(src string, dest string) error {
	// src must be readable
	srcfile, e := os.OpenFile(src, os.O_RDONLY, 0)
	if e != nil {
		return e
	}
	defer srcfile.Close()

	srcstat, e := srcfile.Stat()
	if e != nil {
		return e
	}

	// create dest file
	destfile, e := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, srcstat.Mode())
	if e != nil {
		return e
	}
	defer destfile.Close()
	_, e = io.Copy(destfile, srcfile)
	if e != nil {
		return e
	}
	return nil
}

func CopyDir(src string, dest string) error {
	return copy(src, dest)
}
