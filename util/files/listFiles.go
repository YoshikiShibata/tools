// Copyright © 2016 Yoshiki Shibata. All rights reserved.

package files

import "os"

// ListFiles returns a slice of file name in the specified directory.
// listFiler function will be invoked to determine if a file should be included in the slice.
func ListFiles(dir string, listFilter func(string) bool) ([]string, error) {
	d, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	result := make([]string, 0, len(files))
	for _, file := range files {
		if listFilter(file.Name()) {
			result = append(result, file.Name())
		}
	}
	return result, nil
}
