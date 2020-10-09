package session

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func MakeSessions() error {
	var err error

	// Get recommended shard count from Discord
	if config.Config.System.ShardCount < 1 {
		gatewayBot, err := ChatBackend.RawSession.GatewayBot()
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
		session.Token = ChatBackend.RawSession.Token
		session.UserAgent = ChatBackend.RawSession.UserAgent
		session.ShouldReconnectOnError = ChatBackend.RawSession.ShouldReconnectOnError
		session.Identify.Intents = ChatBackend.RawSession.Identify.Intents
		err = session.Open()
		if err != nil {
			return err
		}
		for _, handler := range ChatBackend.EventHandlers {
			session.AddHandler(handler)
		}
		ChatBackend.ShardSessions = append(ChatBackend.ShardSessions, session)
	}
	return nil
}

func FetchGuildSession(gid string) (*discordgo.Session, error) {
	if !config.Config.System.Shard {
		return ChatBackend.RawSession, nil
	}
	ID, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		return nil, err
	}
	return ChatBackend.ShardSessions[(ID>>22)%int64(config.Config.System.ShardCount)], nil
}
