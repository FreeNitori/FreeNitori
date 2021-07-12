package config

var def = conf{
	System: system{
		LogLevel:       "info",
		LogPath:        "log",
		Socket:         "sock",
		Database:       "db",
		Prefix:         "9 ",
		BackupInterval: 0,
		Administrator:  0,
		Operator:       []int{},
	},
	WebServer: ws{
		Host:                "0.0.0.0",
		Port:                7777,
		BaseURL:             "http://localhost:7777/",
		Unix:                false,
		ForwardedByClientIP: false,
		Secret:              "RANDOM_STRING",
		RateLimit:           1000,
		RateLimitPeriod:     3600,
	},
	Discord: discord{
		Token:           "INSERT_TOKEN_HERE",
		ClientSecret:    "INSERT_CLIENT_SECRET_HERE",
		Presence:        "Manuals: 9 man",
		Shard:           false,
		ShardCount:      8,
		CachePerChannel: 1000,
	},
}
