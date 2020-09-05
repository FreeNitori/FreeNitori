package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"math/rand"
	"strconv"
	"strings"
)

func init() {
	multiplexer.NotTargeted = append(multiplexer.NotTargeted, AdvanceExperience)
	ExperienceCategory.Register(level, "level", []string{"rank", "experience", "exp"}, "Query experience level.")
}

func AdvanceExperience(context *multiplexer.Context) {
	var err error

	// Not do anything if private or bot
	if context.IsPrivate {
		return
	}

	// Also don't do anything if experience system is disabled
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if err != nil {
		return
	}
	if !expEnabled {
		return
	}

	previousExp, err := config.GetMemberExp(context.Author, context.Guild)
	if !context.HandleError(err, config.Debug) {
		return
	}
	advancedExp := previousExp + rand.Intn(10) + 5
	err = config.SetMemberExp(context.Author, context.Guild, advancedExp)
	if !context.HandleError(err, config.Debug) {
		return
	}
	advancedLevel := config.ExpToLevel(advancedExp)
	if advancedLevel > config.ExpToLevel(previousExp) {
		levelupMessage, err := config.GetCustomizableMessage(context.Guild.ID, "levelup")
		if !context.HandleError(err, config.Debug) {
			return
		}
		replacer := strings.NewReplacer("$USER", context.Author.Mention(), "$LEVEL", strconv.Itoa(advancedLevel))
		context.SendMessage(replacer.Replace(levelupMessage))
	}
}

func level(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
	}
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err, config.Debug) {
		return
	}
	if !expEnabled {
		context.SendMessage(state.FeatureDisabled)
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
	if !context.HandleError(err, config.Debug) {
		return
	}
	levelValue := config.ExpToLevel(expValue)
	baseExpValue := config.LevelToExp(levelValue)
	embed.AddField("Level", strconv.Itoa(levelValue), true)
	embed.AddField("Experience", strconv.Itoa(expValue-baseExpValue)+"/"+strconv.Itoa(config.LevelToExp(levelValue+1)-baseExpValue), true)
	embed.SetThumbnail(context.Author.AvatarURL("128"))
	context.SendEmbed(embed)
}
