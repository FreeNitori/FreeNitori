package static

import (
	"net/http"
	"strings"
)

// HTTPFileSystem wraps around http.FileSystem.
type HTTPFileSystem struct {
	filesystem http.FileSystem
}

// Open opens a file.
func (instance *HTTPFileSystem) Open(name string) (http.File, error) {
	return instance.filesystem.Open(name)
}

// Exists returns if a path exists or not.
func (instance *HTTPFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := instance.filesystem.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

// FileSystem returns an HTTPFileSystem.
func FileSystem(fs http.FileSystem) *HTTPFileSystem {
	return &HTTPFileSystem{
		filesystem: fs,
	}
}
