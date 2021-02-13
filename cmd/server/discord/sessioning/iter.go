package sessioning

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

// FetchChannel fetches a channel.
func FetchChannel(guild *discordgo.Guild, channelGeneric, channelID string) *discordgo.Channel {
	for _, channel := range guild.Channels {
		if channelID != "" {
			if channelID == channel.ID {
				return channel
			}
		} else {
			if channelGeneric == channel.Name || channelGeneric == channel.ID || channelGeneric == channel.Mention() {
				return channel
			}
		}
	}
	return nil
}

// FetchGuildSession fetches a session containing a guild from an ID, useful for shard scenario.
func FetchGuildSession(gid string) (*discordgo.Session, error) {
	if !config.Config.Discord.Shard {
		return state.RawSession, nil
	}
	ID, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		return nil, err
	}
	return state.ShardSessions[(ID>>22)%int64(config.Config.Discord.ShardCount)], nil
}

// FetchGuild fetches a guild from an ID.
func FetchGuild(gid string) *discordgo.Guild {
	if _, err := strconv.Atoi(gid); err != nil {
		return nil
	}
	guildSession, err := FetchGuildSession(gid)
	var guild *discordgo.Guild
	if err == nil {
		for _, guildIter := range guildSession.State.Guilds {
			if guildIter.ID == gid {
				guild = guildIter
				break
			}
		}
	}
	return guild
}

// FetchUser fetches a user from an ID.
func FetchUser(uid string) (*discordgo.User, error) {
	if _, err := strconv.Atoi(uid); err != nil {
		return nil, errors.New("invalid snowflake")
	}
	return state.RawSession.User(uid)
}
