package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"strconv"
)

func (*Handlers) Level(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(GuildOnly, "generating guild only message from level query handler")
	}
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if err != nil {
		multiplexer.Logger.Warning(fmt.Sprintf("Failed to obtain experience enabler information, %s", err))
		context.SendMessage(ErrorOccurred, "generating error message")
		return
	}
	if !expEnabled {
		context.SendMessage(FeatureDisabled, "generating disabled message")
		return
	}
	embed := formatter.NewEmbed("Experience Level", context.Author.Username+"#"+context.Author.Discriminator)
	if len(context.Member.Roles) > 0 {
		for _, role := range context.Guild.Roles {
			if role.ID == context.Member.Roles[0] {
				embed.Color = role.Color
				break
			}
		}
	}
	expValue, err := config.GetMemberExp(context.Member)
	if err != nil {
		multiplexer.Logger.Warning(fmt.Sprintf("Failed to obtain experience information, %s", err))
		context.SendMessage(ErrorOccurred, "generating error message")
		return
	}
	levelValue := config.ExpToLevel(expValue)
	embed.AddField("Level", strconv.Itoa(levelValue), true)
	embed.AddField("Experience", strconv.Itoa(expValue)+"/"+strconv.Itoa(config.LevelToExp(expValue+1)), true)
	embed.SetThumbnail(context.Author.AvatarURL("128"))
}
