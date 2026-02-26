package helpers

import (
	"fmt"
	"io/fs"
	"os"
)

func AllUTMLFiles(dirname string) []string {
	root := os.DirFS(dirname)

	utmlFiles, err := fs.Glob(root, "*.utml")
	if err != nil {
		panic(fmt.Sprintf("CANNOT FIND FILES fuck!!!! err=%s", err.Error()))
	}
	return utmlFiles
}
