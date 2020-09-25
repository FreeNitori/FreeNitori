package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
)

func init() {
	multiplexer.NotTargeted = append(multiplexer.NotTargeted, AdvanceExperience)
	ExperienceCategory.Register(level, "level", []string{"rank", "experience", "exp"}, "Query experience level.")
	ExperienceCategory.Register(setrank, "setrank", []string{"rankset"}, "Configure ranked roles.")
}

func AdvanceExperience(context *multiplexer.Context) {
	var err error

	// Not do anything if private
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

	// Obtain experience value of user
	previousExp, err := config.GetMemberExp(context.Author, context.Guild)
	if !context.HandleError(err, config.Debug) {
		return
	}

	// Calculate and set new experience value
	advancedExp := previousExp + rand.Intn(10) + 5
	err = config.SetMemberExp(context.Author, context.Guild, advancedExp)
	if !context.HandleError(err, config.Debug) {
		return
	}

	// Calculate new level value and see if it is advanced as well, and congratulate user if it did
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

	// Doesn't work in private messages
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
	}

	// Checks if feature is enabled
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err, config.Debug) {
		return
	}
	if !expEnabled {
		context.SendMessage(state.FeatureDisabled)
		return
	}

	// Get the member
	var member *discordgo.Member
	if len(context.Fields) > 1 {
		member = context.GetMember(context.StitchFields(1))
	} else {
		member = context.Create.Member
		member.User = context.Author
	}

	// Bail out if nothing is get
	if member == nil {
		context.SendMessage(state.MissingUser)
		return
	}

	// Make the message
	embed := formatter.NewEmbed("Experience Level", member.User.Username+"#"+member.User.Discriminator)
	embed.Color = context.Session.State.UserColor(context.Author.ID, context.Create.ChannelID)
	expValue, err := config.GetMemberExp(member.User, context.Guild)
	if !context.HandleError(err, config.Debug) {
		return
	}
	levelValue := config.ExpToLevel(expValue)
	baseExpValue := config.LevelToExp(levelValue)
	embed.AddField("Level", strconv.Itoa(levelValue), true)
	embed.AddField("Experience", strconv.Itoa(expValue-baseExpValue)+"/"+strconv.Itoa(config.LevelToExp(levelValue+1)-baseExpValue), true)
	embed.SetThumbnail(member.User.AvatarURL("128"))
	context.SendEmbed(embed)
}

func setrank(context *multiplexer.Context) {

	// Doesn't work in private messages
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
	}

	// Checks if feature is enabled
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err, config.Debug) {
		return
	}
	if !expEnabled {
		context.SendMessage(state.FeatureDisabled)
		return
	}

	// Deny access to anyone that does not have permission Administrator
	if !context.HasPermission(discordgo.PermissionAdministrator) {
		context.SendMessage(state.PermissionDenied)
		return
	}

	switch len(context.Fields) {
	case 0:
		embed := formatter.NewEmbed("Ranked Roles", "Configure ranked roles.")
		context.SendEmbed(embed)
	}
}
