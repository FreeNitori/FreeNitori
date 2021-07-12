// +build !windows

package config

var def = conf{
	System: system{
		LogLevel:       "info",
		LogPath:        "/var/log/freenitori",
		Socket:         "/tmp/nitori",
		Database:       "/var/db/freenitori",
		Prefix:         "9 ",
		BackupInterval: 28800,
		Administrator:  0,
		Operator:       []int{},
	},
	WebServer: ws{
		Host:                "127.0.0.1",
		Port:                7777,
		BaseURL:             "http://localhost:7777/",
		Unix:                false,
		ForwardedByClientIP: true,
		Secret:              "RANDOM_STRING",
		RateLimit:           1000,
		RateLimitPeriod:     3600,
	},
	Discord: discord{
		Token:           "INSERT_TOKEN_HERE",
		ClientSecret:    "INSERT_CLIENT_SECRET_HERE",
		Presence:        "Manuals: 9 man",
		Shard:           true,
		ShardCount:      8,
		CachePerChannel: 1024,
	},
}
