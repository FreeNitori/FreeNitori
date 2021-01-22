// Structs used for JSON marshalling by the web server.
package datatypes

import (
	"git.randomchars.net/RandomChars/FreeNitori/binaries/static"
	"net/http"
	"strings"
	"time"
)

// Error messages and stuff
const (
	InternalServerError   = "Internal Server Error"
	NoSuchFileOrDirectory = "No such file or directory"
	ServiceUnavailable    = "Service Unavailable"
	BadRequest            = "Bad Request"
)

type H map[string]interface{}

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

func Public() *fileSystem {
	return &fileSystem{
		filesystem: static.AssetFile(),
	}
}

type LeaderboardEntry struct {
	User       UserInfo
	Experience int
	Level      int
}

type GuildInfo struct {
	Name    string
	ID      string
	IconURL string
	Members []UserInfo
}

type UserInfo struct {
	Name          string
	ID            string
	AvatarURL     string
	Discriminator string
	CreationTime  time.Time
	Bot           bool
}
