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
		context.SendMessage(GuildOnly)
	}
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if err != nil {
		multiplexer.Logger.Warning(fmt.Sprintf("Failed to obtain experience enabler information, %s", err))
		context.SendMessage(ErrorOccurred)
		return
	}
	if !expEnabled {
		context.SendMessage(FeatureDisabled)
		return
	}
	embed := formatter.NewEmbed("Experience Level", context.Author.Username+"#"+context.Author.Discriminator)
	if len(context.Create.Member.Roles) > 0 {
		for _, role := range context.Guild.Roles {
			if role.ID == context.Create.Member.Roles[0] {
				embed.Color = role.Color
				break
			}
		}
	}
	expValue, err := config.GetMemberExp(context.Author, context.Guild)
	if err != nil {
		multiplexer.Logger.Warning(fmt.Sprintf("Failed to obtain experience information, %s", err))
		context.SendMessage(ErrorOccurred)
		return
	}
	levelValue := config.ExpToLevel(expValue)
	baseExpValue := config.LevelToExp(levelValue)
	embed.AddField("Level", strconv.Itoa(levelValue), true)
	embed.AddField("Experience", strconv.Itoa(expValue-baseExpValue)+"/"+strconv.Itoa(config.LevelToExp(levelValue+1)-baseExpValue), true)
	embed.SetThumbnail(context.Author.AvatarURL("128"))
	context.SendEmbed(embed)
}
