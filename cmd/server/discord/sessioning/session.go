package sessioning

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
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
		session.State.MaxMessageCount = state.RawSession.State.MaxMessageCount
		err = session.Open()
		if err != nil {
			return err
		}
		state.Multiplexer.SessionRegisterHandlers(session)
		log.Infof("Shard %s ready.", strconv.Itoa(i))
		state.ShardSessions = append(state.ShardSessions, session)
	}
	return nil
}
