// Structs used for JSON marshalling by the web server.
package jsontypes

import "time"

type GuildInfo struct {
	Name    string
	ID      string
	IconURL string
	Members []*UserInfo
}

type UserInfo struct {
	Name          string
	ID            string
	AvatarURL     string
	Discriminator string
	CreationTime  time.Time
	Bot           bool
}
