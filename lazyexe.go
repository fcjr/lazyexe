package lazyexe

import (
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

var prefix = "lazyexe-"

type LazyExe struct {
	mu      sync.Mutex
	bytes   []byte
	tmpFile *os.File // non nil once exe is loaded
}

// New creates a new lazyexe file which will be written to disk when the path is initially requested.
// supsequent requests for the path will return the same instance of the executable on disk.
// the caller is responsible for calling Cleanup() when the lazyexe is no longer needed.
func New(b []byte) *LazyExe {
	le := &LazyExe{
		bytes: b,
	}
	return le
}

// Path returns the temporary path to the requested executable written on disk
// lazily writing and chmoding it when first requested.
func (le *LazyExe) Path() (string, error) {
	le.mu.Lock()
	defer le.mu.Unlock()
	if le.tmpFile != nil {
		return le.tmpFile.Name(), nil
	}

	var err error
	if runtime.GOOS == "windows" {
		prefix += "*.exe"
	}
	le.tmpFile, err = ioutil.TempFile(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}

	if _, err = le.tmpFile.Write(le.bytes); err != nil {
		return "", err
	}

	if err := le.tmpFile.Chmod(0700); err != nil {
		return "", err
	}

	if err := le.tmpFile.Close(); err != nil {
		return "", err
	}

	return le.tmpFile.Name(), nil
}

// Cleanup removes any temporary files LazyExe has written to disk.
// It is expected to be called when the lazyexe is no longer needed.
func (le *LazyExe) Cleanup() error {
	le.mu.Lock()
	defer le.mu.Unlock()
	if le.tmpFile != nil {
		if err := os.Remove(le.tmpFile.Name()); err != nil {
			return err
		}
		le.tmpFile = nil
	}
	return nil
}
