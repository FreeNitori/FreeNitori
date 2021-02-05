package internals

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
)

func init() {
	multiplexer.NotTargeted = append(multiplexer.NotTargeted, AdvanceExperience)
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "level",
		AliasPatterns: []string{"rank", "experience", "exp"},
		Description:   "Query experience level.",
		Category:      multiplexer.ExperienceCategory,
		Handler:       level,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "setrank",
		AliasPatterns: []string{"rankset"},
		Description:   "Configure ranked roles.",
		Category:      multiplexer.ExperienceCategory,
		Handler:       setrank,
	})
	overrides.RegisterSimpleEntry(overrides.SimpleConfigurationEntry{
		Name:         "experience",
		FriendlyName: "Chat Experience System",
		Description:  "Toggle chat experience system.",
		DatabaseKey:  "exp_enable",
		Cleanup:      func(context *multiplexer.Context) {},
		Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
			if *input != "toggle" {
				return false, true
			}
			pre, err := config.ExpEnabled(context.Guild.ID)
			if !context.HandleError(err) {
				return true, false
			}
			switch pre {
			case true:
				*input = "false"
			case false:
				*input = "true"
			}
			return true, true
		},
		Format: func(context *multiplexer.Context, value string) (string, string, bool) {
			pre, err := config.ExpEnabled(context.Guild.ID)
			if !context.HandleError(err) {
				return "", "", false
			}
			description := fmt.Sprintf("Toggle by issuing command `%sconf experience toggle`.", context.Prefix())
			switch pre {
			case true:
				return "Chat experience system enabled", description, true
			case false:
				return "Chat experience system disabled", description, true
			}
			return "", "", false
		},
	})
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
	if !context.HandleError(err) {
		return
	}

	// Calculate and set new experience value
	advancedExp := previousExp + rand.Intn(10) + 5
	err = config.SetMemberExp(context.Author, context.Guild, advancedExp)
	if !context.HandleError(err) {
		return
	}

	// Calculate new level value and see if it is advanced as well, and congratulate user if it did
	advancedLevel := config.ExpToLevel(advancedExp)
	if advancedLevel > config.ExpToLevel(previousExp) {
		levelupMessage, err := config.GetCustomizableMessage(context.Guild.ID, "levelup")
		if !context.HandleError(err) {
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
		return
	}

	// Checks if feature is enabled
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err) {
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
	embed := embedutil.NewEmbed("Experience Level", member.User.Username+"#"+member.User.Discriminator)
	embed.Color = context.Session.State.UserColor(context.Author.ID, context.Create.ChannelID)
	expValue, err := config.GetMemberExp(member.User, context.Guild)
	if !context.HandleError(err) {
		return
	}
	levelValue := config.ExpToLevel(expValue)
	baseExpValue := config.LevelToExp(levelValue)
	embed.AddField("Level", strconv.Itoa(levelValue), true)
	embed.AddField("Experience", strconv.Itoa(expValue-baseExpValue)+"/"+strconv.Itoa(config.LevelToExp(levelValue+1)-baseExpValue), true)
	embed.SetThumbnail(member.User.AvatarURL("128"))
	context.SendEmbed("", embed)
}

func setrank(context *multiplexer.Context) {

	// Doesn't work in private messages
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
	}

	// Checks if feature is enabled
	expEnabled, err := config.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err) {
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
		embed := embedutil.NewEmbed("Ranked Roles", "Configure ranked roles.")
		context.SendEmbed("", embed)
	}
}
