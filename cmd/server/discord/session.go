package discord

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

// MakeSessions opens all sessions of shards.
func MakeSessions() error {
	var err error

	// Get recommended shard count from Discord
	if config.Config.Discord.ShardCount < 1 {
		gatewayBot, err := state.RawSession.GatewayBot()
		if err != nil {
			return err
		}
		config.Config.Discord.ShardCount = gatewayBot.Shards
	}

	// Make sure it doesn't end up being 0 shards
	if config.Config.Discord.ShardCount == 0 {
		config.Config.Discord.ShardCount = 1
	}

	// Make the sessions
	for i := 0; i < config.Config.Discord.ShardCount; i++ {
		time.Sleep(time.Millisecond * 100)
		session, _ := discordgo.New()
		session.ShardCount = config.Config.Discord.ShardCount
		session.ShardID = i
		session.Token = state.RawSession.Token
		session.UserAgent = state.RawSession.UserAgent
		session.ShouldReconnectOnError = state.RawSession.ShouldReconnectOnError
		session.Identify.Intents = state.RawSession.Identify.Intents
		err = session.Open()
		if err != nil {
			return err
		}
		for _, handler := range multiplexer.EventHandlers {
			session.AddHandler(handler)
		}
		log.Infof("Shard %s ready.", strconv.Itoa(i))
		state.ShardSessions = append(state.ShardSessions, session)
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
