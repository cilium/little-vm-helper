package kernels

import (
	"fmt"
	"os"
)

// directoryExists returns:
//   true,  nil: if a directory with name dir exists
//   false, nil: if a directory with name dir does not exist
//   false, err: if somethign unexpected happened
func directoryExists(dir string) (bool, error) {
	st, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err == nil {
		if st.IsDir() {
			return true, nil
		}
		return false, fmt.Errorf("`%s` exists, but is not a directory", dir)
	}

	return false, fmt.Errorf("error accessing `%s`: %w", dir, err)
}

func regularFileExists(fname string) (bool, error) {
	st, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err == nil {
		if st.Mode().IsRegular() {
			return true, nil
		}
		return false, fmt.Errorf("`%s` exists, but is not a regular file", fname)
	}

	return false, fmt.Errorf("error accessing `%s`: %w", fname, err)
}
