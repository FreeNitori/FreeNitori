package routes

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/db"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

func init() {
	state.Multiplexer.NotTargeted = append(state.Multiplexer.NotTargeted, AdvanceExperience)
	state.Multiplexer.GuildMemberAdd = append(state.Multiplexer.GuildMemberAdd, memberAddRank)
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "level",
		AliasPatterns: []string{"rank", "experience", "exp"},
		Description:   "Query experience level.",
		Category:      multiplexer.ExperienceCategory,
		Handler:       level,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "leaderboard",
		AliasPatterns: []string{"lb"},
		Description:   "Display URL of leaderboard.",
		Category:      multiplexer.ExperienceCategory,
		Handler:       leaderboard,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "exp2level",
		AliasPatterns: []string{},
		Description:   "",
		Category:      multiplexer.ExperienceCategory,
		Handler:       exp2level,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "level2exp",
		AliasPatterns: []string{},
		Description:   "",
		Category:      multiplexer.ExperienceCategory,
		Handler:       level2exp,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "setexp",
		AliasPatterns: []string{},
		Description:   "",
		Category:      multiplexer.ExperienceCategory,
		Handler:       setexp,
	})
	overrides.RegisterComplexEntry(overrides.ComplexConfigurationEntry{
		Name:         "experience",
		FriendlyName: "Chat Experience System",
		Description:  "Configure chat experience related options.",
		Entries: []overrides.SimpleConfigurationEntry{
			{
				Name:         "enable",
				FriendlyName: "Enable Experience System",
				Description:  "Toggle chat experience system.",
				DatabaseKey:  "exp_enable",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if *input != "toggle" {
						return false, true
					}
					pre, err := db.ExpEnabled(context.Guild.ID)
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
					pre, err := db.ExpEnabled(context.Guild.ID)
					if !context.HandleError(err) {
						return "", "", false
					}
					description := fmt.Sprintf("Toggle by issuing command `%sconf experience enable toggle`.", context.Prefix())
					switch pre {
					case true:
						return "Chat experience system enabled", description, true
					case false:
						return "Chat experience system disabled", description, true
					}
					return "", "", false
				},
			},
		},
		CustomEntries: []overrides.CustomConfigurationEntry{
			{
				Name:        "rank",
				Description: "Configure ranked roles.",
				Handler: func(context *multiplexer.Context) {

					// Checks if feature is enabled
					expEnabled, err := db.ExpEnabled(context.Guild.ID)
					if !context.HandleError(err) {
						return
					}
					if !expEnabled {
						context.SendMessage(multiplexer.FeatureDisabled)
						return
					}

					switch len(context.Fields) {
					case 5:
						level, err := strconv.Atoi(context.Fields[3])
						if err != nil {
							context.SendMessage(multiplexer.InvalidArgument)
							return
						}
						if level < 0 {
							context.SendMessage(multiplexer.InvalidArgument)
							return
						}
						bindings, err := db.GetRankBinds(context.Guild)
						if !context.HandleError(err) {
							return
						}
						if len(bindings) > 16 {
							context.SendMessage(multiplexer.InvalidArgument)
							return
						}
						for _, role := range context.Guild.Roles {
							if context.Fields[4] == role.ID || context.Fields[4] == role.Name || context.Fields[4] == role.Mention() {
								if role.Managed {
									context.SendMessage(multiplexer.InvalidArgument)
									return
								}
								err = db.SetRankBind(context.Guild, level, role)
								if !context.HandleError(err) {
									return
								}
								context.SendMessage(fmt.Sprintf("Successfully bound level %v to role %s.", level, role.Mention()))
								return
							}
						}
						context.SendMessage(multiplexer.InvalidArgument)
					case 4:
						level, err := strconv.Atoi(context.Fields[3])
						if err != nil {
							context.SendMessage(multiplexer.InvalidArgument)
							return
						}
						binding, err := db.GetRankBind(context.Guild, level)
						if !context.HandleError(err) {
							return
						}
						if binding == "" {
							context.SendMessage(multiplexer.InvalidArgument)
							return
						}
						err = db.UnsetRankBind(context.Guild, strconv.Itoa(level))
						if !context.HandleError(err) {
							return
						}
						context.SendMessage(fmt.Sprintf("Successfully removed ranked role binding on level %v.", level))
					case 3:
						bindings, err := db.GetRankBinds(context.Guild)
						if !context.HandleError(err) {
							return
						}
						var embed embedutil.Embed
						if len(bindings) == 0 {
							embed = embedutil.New("Ranked Roles", "No ranked roles are set.")
							embed.Color = multiplexer.KappaColor
						} else {
							embed = embedutil.New("Ranked Roles", "")
							embed.Color = multiplexer.KappaColor
							for level, roleID := range bindings {
								var role *discordgo.Role
								for _, r := range context.Guild.Roles {
									if r.ID == roleID {
										role = r
										break
									}
								}
								if role == nil {
									err = db.UnsetRankBind(context.Guild, level)
									if !context.HandleError(err) {
										return
									}
									continue
								}
								embed.AddField("Level "+level, role.Mention(), false)
							}
						}
						context.SendEmbed("", embed)
					default:
						context.SendMessage(multiplexer.InvalidArgument)
						return
					}
				},
			},
		},
	})
}

func memberAddRank(context *multiplexer.Context) {
	if context.User.Bot {
		return
	}

	if len(context.Guild.Channels) == 0 {
		return
	}

	// If Nitori has permission
	permissions, err := context.Session.State.UserChannelPermissions(context.Session.State.User.ID, context.Guild.Channels[0].ID)
	if !(err == nil && (permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles)) {
		return
	}

	exp, err := db.GetMemberExp(context.User, context.Guild)
	if err != nil {
		return
	}
	bindings, err := db.GetRankBinds(context.Guild)
	if err != nil {
		return
	}
	memberLevel := ExpToLevel(exp)
	for i := 0; i <= memberLevel; i++ {
		roleID := bindings[strconv.Itoa(i)]
		if roleID != "" {
			for _, r := range context.Guild.Roles {
				if r.ID == roleID {
					err = context.Session.GuildMemberRoleAdd(context.Guild.ID, context.User.ID, roleID)
					if err != nil {
						return
					}
				}
			}
		}
	}
}

// AdvanceExperience advances experience of author.
func AdvanceExperience(context *multiplexer.Context) {
	var err error

	// Not do anything if private
	if context.IsPrivate {
		return
	}

	// Check if Channel is EXPBL
	if context.Channel != nil {
		if strings.HasSuffix(context.Channel.Topic, "[EXPBL]") {
			return
		}
	} else {
		return
	}

	// Also don't do anything if experience system is disabled
	expEnabled, err := db.ExpEnabled(context.Guild.ID)
	if err != nil {
		return
	}
	if !expEnabled {
		return
	}

	// Obtain experience value of user
	previousExp, err := db.GetMemberExp(context.User, context.Guild)
	if !context.HandleError(err) {
		return
	}

	// Calculate and set new experience value
	advancedExp := previousExp + rand.Intn(10) + 5
	err = db.SetMemberExp(context.User, context.Guild, advancedExp)
	if !context.HandleError(err) {
		return
	}

	// Calculate new level value and see if it is advanced as well, and congratulate user if it did
	advancedLevel := ExpToLevel(advancedExp)
	if advancedLevel > ExpToLevel(previousExp) {
		// Do role assignment if Nitori has permission
		permissions, err := context.Session.State.UserChannelPermissions(context.Session.State.User.ID, context.Message.ChannelID)
		if err == nil && (permissions&discordgo.PermissionManageRoles == discordgo.PermissionManageRoles) {
			bindings, err := db.GetRankBinds(context.Guild)
			if !context.HandleError(err) {
				return
			}
			if bindings[strconv.Itoa(advancedLevel)] != "" {
				for _, r := range context.Guild.Roles {
					if bindings[strconv.Itoa(advancedLevel)] == r.ID {
						err = context.Session.GuildMemberRoleAdd(context.Guild.ID, context.User.ID, r.ID)
						if !context.HandleError(err) {
							return
						}
					}
				}
			}
		}
		levelupMessage, err := db.GetCustomizableMessage(context.Guild.ID, "levelup")
		if !context.HandleError(err) {
			return
		}
		replacer := strings.NewReplacer("$USER", context.User.Mention(), "$LEVEL", strconv.Itoa(advancedLevel))
		context.SendMessage(replacer.Replace(levelupMessage))
	}
}

func level(context *multiplexer.Context) {

	// Doesn't work in private messages
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
		return
	}

	// Checks if feature is enabled
	expEnabled, err := db.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err) {
		return
	}
	if !expEnabled {
		context.SendMessage(multiplexer.FeatureDisabled)
		return
	}

	// Get the member
	var member *discordgo.Member
	if len(context.Fields) > 1 {
		member = context.GetMember(context.StitchFields(1))
	} else {
		member = context.Member
		member.User = context.User
	}

	// Bail out if nothing is get
	if member == nil {
		context.SendMessage(multiplexer.MissingUser)
		return
	}

	// Make the message
	embed := embedutil.New("Experience Level", member.User.Username+"#"+member.User.Discriminator)
	embed.Color = context.Session.State.UserColor(context.User.ID, context.Channel.ID)
	expValue, err := db.GetMemberExp(member.User, context.Guild)
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
		context.SendMessage(multiplexer.GuildOnly)
		return
	}
	enabled, err := db.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err) {
		return
	}
	if !enabled {
		context.SendMessage(multiplexer.FeatureDisabled)
		return
	}
	embed := embedutil.New("Leaderboard",
		fmt.Sprintf("Click [here](%sleaderboard.html#%s) to view the leaderboard.",
			config.WebServer.BaseURL,
			context.Guild.ID))
	embed.Color = multiplexer.KappaColor
	context.SendEmbed("", embed)
}

func level2exp(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}
	if len(context.Fields) == 2 {
		lvl, err := strconv.Atoi(context.Fields[1])
		if err != nil {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		context.SendMessage(fmt.Sprintf("%v levels is %v exp.", lvl, LevelToExp(lvl)))
	} else {
		context.SendMessage(multiplexer.InvalidArgument)
	}
}

func exp2level(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}
	if len(context.Fields) == 2 {
		exp, err := strconv.Atoi(context.Fields[1])
		if err != nil {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		context.SendMessage(fmt.Sprintf("%v exp is %v levels.", exp, ExpToLevel(exp)))
	} else {
		context.SendMessage(multiplexer.InvalidArgument)
	}
}

func setexp(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
		return
	}
	enabled, err := db.ExpEnabled(context.Guild.ID)
	if !context.HandleError(err) {
		return
	}
	if !enabled {
		context.SendMessage(multiplexer.FeatureDisabled)
		return
	}
	if len(context.Fields) == 3 {
		exp, err := strconv.Atoi(context.Fields[1])
		if err != nil {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		member := context.GetMember(context.Fields[2])
		if member == nil {
			context.SendMessage(multiplexer.MissingUser)
			return
		}
		err = db.SetMemberExp(member.User, context.Guild, exp)
		if !context.HandleError(err) {
			return
		}
		context.SendMessage(fmt.Sprintf("Successfully set experience of %s to %v.", member.User.Username, exp))
	} else {
		context.SendMessage(multiplexer.InvalidArgument)
	}
}

// LevelToExp calculates amount of experience from a level integer.
func LevelToExp(level int) int {
	return int(1000.0 * (math.Pow(float64(level), 1.25)))
}

// ExpToLevel calculates amount of levels from an experience integer.
func ExpToLevel(exp int) int {
	return int(math.Pow(float64(exp)/1000, 1.0/1.25))
}
