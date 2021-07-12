package config

type conf struct {
	System    system
	WebServer ws
	Discord   discord
}

type system struct {
	LogLevel       string
	LogPath        string
	Socket         string
	Database       string
	Prefix         string
	BackupInterval int
	Administrator  int
	Operator       []int
}

type ws struct {
	Host                string
	Port                int
	BaseURL             string
	Unix                bool
	ForwardedByClientIP bool
	Secret              string
	RateLimit           int
	RateLimitPeriod     int
}

type discord struct {
	Token           string
	ClientSecret    string
	Presence        string
	Shard           bool
	ShardCount      int
	CachePerChannel int
}
