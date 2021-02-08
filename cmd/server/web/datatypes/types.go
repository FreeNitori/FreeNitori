// Package datatypes primarily contains structs used for JSON operations.
package datatypes

import "time"

// Error messages.
const (
	InternalServerError   = "Internal Server Error"
	NoSuchFileOrDirectory = "No such file or directory"
	ServiceUnavailable    = "Service Unavailable"
	BadRequest            = "Bad Request"
)

// H is a shortcut to a string to interface map.
type H map[string]interface{}

// LeaderboardEntry represents an entry in a leaderboard.
type LeaderboardEntry struct {
	User       UserInfo
	Experience int
	Level      int
}

// GuildInfo represents information of a guild.
type GuildInfo struct {
	Name    string
	ID      string
	IconURL string
	Members []UserInfo
}

// UserInfo represents information of a user.
type UserInfo struct {
	Name          string
	ID            string
	AvatarURL     string
	Discriminator string
	CreationTime  time.Time
	Bot           bool
}
