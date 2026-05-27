package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4" // for being able to use the "/**/*.<extension>" syntax in GLOB syntax
)

func AllUTMLFilesUNSAFE(dirname string) []string {
	r, err := AllUTMLFiles(dirname)
	if err != nil {
		panic(err.Error())
	}
	return r
}

func AllUTMLFiles(path string) ([]string, error) {
	return AllFiles(path, "**/*.utml", "**/*.json")
}

const DATASET_DIR = "../DATASETS"

var DATASET_FILE_GLOBS = []string{
	"**/q/**/*.json",
	"**/q/**/*.utml",
}

func AllDatasetFiles() ([]string, error) {
	return AllFiles(DATASET_DIR, DATASET_FILE_GLOBS...)
}

// Supports the double-star syntax in GLOB ("**/*.go" for example)
func AllFiles(path string, globs ...string) ([]string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch UTML files: %s", err.Error())
	}

	if !stat.IsDir() {
		return []string{path}, nil
	}

	root := os.DirFS(path)

	files, err := []string{}, nil
	for _, gl := range globs {
		var fs []string
		fs, err = doublestar.Glob(root, gl) //directly assign to the outer 'err'
		if err != nil {
			break
		}
		files = append(files, fs...)
	}

	if err != nil {
		return nil, fmt.Errorf("Cannot find UTML files in '%s': %s", path, err.Error())
	}

	for i, file := range files {
		files[i] = filepath.Join(path, file)
	}

	return files, nil
}
