package script

import (
	"io/ioutil"
	"os"
	"testing"
)

// MkTmpDir creates a temporary directory usiing ioutil.TempDir and calls t.Fatal if the attempt fails. On success, it
// returns:
// - the created directory
// - a no-arg function which deletes the temp directory and os.Chdir to the current working directory
func MkTmpDir(t *testing.T, errid string) (string, func()) {
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(errid, err)
	}
	tempDir, err := ioutil.TempDir("", errid)
	if err != nil {
		t.Fatal(errid, err)
	}

	return tempDir, func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatal(errid, err)
		}
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatal(errid, err)
		}
	}
}
