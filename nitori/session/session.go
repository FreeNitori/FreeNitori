package session

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func MakeSessions() error {
	var err error

	// Get recommended shard count from Discord
	if config.Config.System.ShardCount < 1 {
		gatewayBot, err := state.RawSession.GatewayBot()
		if err != nil {
			return err
		}
		config.Config.System.ShardCount = gatewayBot.Shards
	}

	// Make sure it doesn't end up being 0 shards
	if config.Config.System.ShardCount == 0 {
		config.Config.System.ShardCount = 1
	}

	// Make the sessions
	for i := 0; i < config.Config.System.ShardCount; i++ {
		time.Sleep(time.Millisecond * 100)
		session, _ := discordgo.New()
		session.ShardCount = config.Config.System.ShardCount
		session.ShardID = i
		session.Token = state.RawSession.Token
		session.UserAgent = state.RawSession.UserAgent
		session.ShouldReconnectOnError = state.RawSession.ShouldReconnectOnError
		session.Identify.Intents = state.RawSession.Identify.Intents
		err = session.Open()
		if err != nil {
			return err
		}
		for _, handler := range state.EventHandlers {
			session.AddHandler(handler)
		}
		state.ShardSessions = append(state.ShardSessions, session)
	}
	return nil
}

func FetchGuildSession(gid string) (*discordgo.Session, error) {
	if !config.Config.System.Shard {
		return state.RawSession, nil
	}
	ID, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		return nil, err
	}
	return state.ShardSessions[(ID>>22)%int64(config.Config.System.ShardCount)], nil
}
