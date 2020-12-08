package jsontypes

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
	Bot           bool
}