package internals

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"math"
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
		Pattern:       "leaderboard",
		AliasPatterns: []string{"lb"},
		Description:   "Display URL of leaderboard.",
		Category:      multiplexer.ExperienceCategory,
		Handler:       leaderboard,
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
	advancedLevel := ExpToLevel(advancedExp)
	if advancedLevel > ExpToLevel(previousExp) {
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
	levelValue := ExpToLevel(expValue)
	baseExpValue := LevelToExp(levelValue)
	embed.AddField("Level", strconv.Itoa(levelValue), true)
	embed.AddField("Experience", strconv.Itoa(expValue-baseExpValue)+"/"+strconv.Itoa(LevelToExp(levelValue+1)-baseExpValue), true)
	embed.SetThumbnail(member.User.AvatarURL("128"))
	context.SendEmbed("", embed)
}

func leaderboard(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
		return
	}
	enabled, err := config.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err) {
		return
	}
	if !enabled {
		context.SendMessage(state.FeatureDisabled)
		return
	}
	embed := embedutil.NewEmbed("Leaderboard",
		fmt.Sprintf("Click [here](%sleaderboard.html#%s) to view the leaderboard.",
			config.Config.WebServer.BaseURL,
			context.Guild.ID))
	embed.Color = state.KappaColor
	context.SendEmbed("", embed)
}

// ExpToLevel calculates amount of experience from a level integer.
func LevelToExp(level int) int {
	return int(1000.0 * (math.Pow(float64(level), 1.25)))
}

// ExpToLevel calculates amount of levels from an experience integer.
func ExpToLevel(exp int) int {
	return int(math.Pow(float64(exp)/1000, 1.0/1.25))
}
