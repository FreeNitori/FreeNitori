package static

import (
	"net/http"
	"strings"
)

type fileSystem struct {
	filesystem http.FileSystem
}

func (instance *fileSystem) Open(name string) (http.File, error) {
	return instance.filesystem.Open(name)
}

func (instance *fileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := instance.filesystem.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func FileSystem(fs http.FileSystem) *fileSystem {
	return &fileSystem{
		filesystem: fs,
	}
}
