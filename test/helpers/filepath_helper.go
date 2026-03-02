package helpers

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

func AllUTMLFiles(dirname string) []string {
	root := os.DirFS(dirname)

	utmlFiles, err := fs.Glob(root, "*.utml")
	if err != nil {
		panic(fmt.Sprintf("CANNOT FIND FILES fuck!!!! err=%s", err.Error()))
	}

	for i, file := range utmlFiles {
		utmlFiles[i] = path.Join(dirname, file)
	}

	return utmlFiles
}
