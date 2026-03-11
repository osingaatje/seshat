package helper

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

func AllUTMLFiles(dirOrFileName string) ([]string, error) {
	stat, err := os.Stat(dirOrFileName)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch UTML files: %s", err.Error())
	}

	if !stat.IsDir() {
		return []string{dirOrFileName}, nil
	}

	root := os.DirFS(dirOrFileName)

	utmlFiles, err := fs.Glob(root, "*.utml")
	if err != nil {
		return nil, fmt.Errorf("Cannot find UTML files in '%s': %s", dirOrFileName, err.Error())
	}

	for i, file := range utmlFiles {
		utmlFiles[i] = path.Join(dirOrFileName, file)
	}

	return utmlFiles, nil
}

func AllUTMLFilesUNSAFE(dirname string) []string {
	r, err := AllUTMLFiles(dirname)
	if err != nil {
		panic(err.Error())
	}
	return r
}
