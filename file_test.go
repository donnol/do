package do

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMkdirAllIfNotExist(t *testing.T) {
	path1 := "./testdir/"
	path2 := "testdir2"
	dir := filepath.Join(path1, path2)
	if err := MkdirAllIfNotExist(dir); err != nil {
		t.Error(err)
	}

	fi, err := os.Stat(dir)
	if err != nil {
		t.Error(err)
	}
	if !fi.IsDir() {
		t.Errorf("bad case %s is not a directory", fi.Name())
	}
	if fi.Name() != path2 {
		t.Errorf("bad case: %s != %s", fi.Name(), path2)
	}

	if err := os.RemoveAll(path1); err != nil {
		t.Error(err)
	}
}
