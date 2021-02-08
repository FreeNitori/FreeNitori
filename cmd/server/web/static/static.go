package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
)

const index = "index.html"

// ServeFileSystem represents a static filesystem that can be served.
type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

// LocalFileSystem represents a local file system.
type LocalFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

// LocalFile returns an instance of a LocalFileSystem serving from a specified directory.
func LocalFile(root string, indexes bool) *LocalFileSystem {
	return &LocalFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

// Exists returns if a path exists.
func (l *LocalFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
				index := path.Join(name, index)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

// ServeRoot returns a middleware handler that serves from a filesystem directory.
func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Serve returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileServer := http.FileServer(fs)
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
