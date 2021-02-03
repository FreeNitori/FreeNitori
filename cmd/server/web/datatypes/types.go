// Structs used for JSON marshalling by the web server.
package datatypes

import "time"

// Error messages and stuff
const (
	InternalServerError   = "Internal Server Error"
	NoSuchFileOrDirectory = "No such file or directory"
	ServiceUnavailable    = "Service Unavailable"
	BadRequest            = "Bad Request"
)

type H map[string]interface{}

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
