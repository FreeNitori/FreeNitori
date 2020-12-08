package discord

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func MakeSessions() error {
	var err error

	// Get recommended shard count from Discord
	if config.Config.Discord.ShardCount < 1 {
		gatewayBot, err := vars.RawSession.GatewayBot()
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
		session.Token = vars.RawSession.Token
		session.UserAgent = vars.RawSession.UserAgent
		session.ShouldReconnectOnError = vars.RawSession.ShouldReconnectOnError
		session.Identify.Intents = vars.RawSession.Identify.Intents
		err = session.Open()
		if err != nil {
			return err
		}
		for _, handler := range vars.EventHandlers {
			session.AddHandler(handler)
		}
		log.Infof("Shard %s ready.", strconv.Itoa(i))
		vars.ShardSessions = append(vars.ShardSessions, session)
	}
	return nil
}

func FetchGuildSession(gid string) (*discordgo.Session, error) {
	if !config.Config.Discord.Shard {
		return vars.RawSession, nil
	}
	ID, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		return nil, err
	}
	return vars.ShardSessions[(ID>>22)%int64(config.Config.Discord.ShardCount)], nil
}

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
