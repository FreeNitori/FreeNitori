// +build !windows

package config

var confDefault = Conf{
	System: struct {
		LogLevel       string
		LogPath        string
		Socket         string
		Database       string
		Prefix         string
		BackupInterval int
		Administrator  int
		Operator       []int
	}{
		LogLevel:       "info",
		LogPath:        "/var/log/freenitori",
		Socket:         "/tmp/nitori",
		Database:       "/var/db/freenitori",
		Prefix:         "env ",
		BackupInterval: 28800,
		Administrator:  0,
		Operator:       []int{},
	},
	WebServer: struct {
		Host                string
		Port                int
		BaseURL             string
		Unix                bool
		ForwardedByClientIP bool
		Secret              string
		RateLimit           int
		RateLimitPeriod     int
	}{
		Host:                "127.0.0.1",
		Port:                7777,
		BaseURL:             "http://localhost:7777/",
		Unix:                false,
		ForwardedByClientIP: true,
		Secret:              "RANDOM_STRING",
		RateLimit:           1000,
		RateLimitPeriod:     3600,
	},
	Discord: struct {
		Token           string
		ClientSecret    string
		Presence        string
		Shard           bool
		ShardCount      int
		CachePerChannel int
	}{
		Token:           "INSERT_TOKEN_HERE",
		ClientSecret:    "INSERT_CLIENT_SECRET_HERE",
		Presence:        "Manuals: env man",
		Shard:           true,
		ShardCount:      8,
		CachePerChannel: 1000,
	},
}
