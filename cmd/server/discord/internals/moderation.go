package internals

import (
	"encoding/json"
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/db"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

type MemberWarning struct {
	Text      string    `json:"text"`
	Time      time.Time `json:"time"`
	GuildID   string    `json:"guild_id"`
	ChannelID string    `json:"channel_id"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
}

func init() {
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "userinfo",
		AliasPatterns: []string{"whois", "lookup", "pfp"},
		Description:   "Lookup a user's detailed information by username, nickname or snowflake.",
		Category:      multiplexer.ModerationCategory,
		Handler:       userinfo,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "guildinfo",
		AliasPatterns: []string{"pfp"},
		Description:   "Lookup a guild's detailed information by snowflake.",
		Category:      multiplexer.ModerationCategory,
		Handler:       guildinfo,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "warn",
		AliasPatterns: []string{"warnings", "warning"},
		Description:   "Lookup warnings associated with a user or assign/clear warning.",
		Category:      multiplexer.ModerationCategory,
		Handler:       warn,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "ban",
		AliasPatterns: []string{""},
		Description:   "Ban a user from the guild",
		Category:      multiplexer.ModerationCategory,
		Handler:       ban,
	})
}

func userinfo(context *multiplexer.Context) {
	var user *discordgo.User
	var member *discordgo.Member

	// Just use the author if there's no arguments
	if len(context.Fields) == 1 {
		user = context.User
		if !context.IsPrivate {
			member = context.Member
		}
	} else {
		argument := context.StitchFields(1)
		member = context.GetMember(argument)
		if member != nil {
			user = member.User
		} else {
			user, _ = context.Session.User(argument)
		}
	}

	// Fail if unable to make sense of the arguments passed
	if user == nil {
		context.SendMessage(multiplexer.MissingUser)
		return
	}

	// Only generate the pfp stuff if that's what's required
	if context.Fields[0] == "pfp" {
		embed := embedutil.New("", "")
		embed.Color = context.Session.State.UserColor(user.ID, context.Channel.ID)
		embed.SetAuthor(user.Username, user.AvatarURL("128"))
		embed.SetImage(user.AvatarURL("4096"))
		context.SendEmbed("", embed)
		return
	}

	// Make the message
	userID, err := strconv.Atoi(user.ID)
	if !context.HandleError(err) {
		return
	}
	embed := embedutil.New("User Information", "")
	embed.Color = context.Session.State.UserColor(user.ID, context.Channel.ID)
	embed.SetThumbnail(user.AvatarURL("1024"))
	embed.AddField("Username", user.Username+"#"+user.Discriminator, member != nil)
	if member != nil {
		var roles string
		if len(member.Roles) == 0 {
			roles = "No Roles"
		} else {
			for _, i := range context.Guild.Roles {
				for _, j := range member.Roles {
					if j == i.ID {
						roles += i.Mention() + "\n"
					}
				}
			}
		}
		if member.Nick != "" {
			embed.AddField("Nickname", member.Nick, true)
		}
		embed.AddField("Roles", roles, false)
	}
	embed.AddField("Registration Date", time.Unix(int64(((userID>>22)+1420070400000)/1000), 0).UTC().Format("Mon, 02 Jan 2006 15:04:05"), true)
	if member != nil {
		joinTime, err := member.JoinedAt.Parse()
		if !context.HandleError(err) {
			return
		}
		embed.AddField("Join Date", joinTime.Format("Mon, 02 Jan 2006 15:04:05"), true)
	}
	embed.SetFooter("ID: " + user.ID)
	context.SendEmbed("", embed)
}

func guildinfo(context *multiplexer.Context) {
	// Guild only
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
		return
	}
	guildID, err := strconv.Atoi(context.Guild.ID)
	if !context.HandleError(err) {
		return
	}
	embed := embedutil.New("Guild Information", "")
	embed.Color = multiplexer.KappaColor
	embed.SetThumbnail(context.Guild.IconURL())
	embed.AddField("Guild Name", context.Guild.Name, true)
	embed.AddField("Member Count", strconv.Itoa(context.Guild.MemberCount), true)
	if len(context.Guild.Roles) > 0 {
		var roles string
		for i := 0; i < len(context.Guild.Roles); i++ {
			roles += context.Guild.Roles[i].Mention() + "\n"
		}
		embed.AddField("Roles", roles, false)
	}
	embed.AddField("Region", context.Guild.Region, true)
	embed.AddField("Locale", context.Guild.PreferredLocale, true)
	embed.AddField("Creation Date", time.Unix(int64(((guildID>>22)+1420070400000)/1000), 0).UTC().Format("Mon, 02 Jan 2006 15:04:05"), true)
	embed.SetFooter("ID: " + context.Guild.ID)
	context.SendEmbed("", embed)
}

func warn(context *multiplexer.Context) {
	// Guild only
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
		return
	}

	// Has permission
	if !context.HasPermission(discordgo.PermissionBanMembers) {
		context.SendMessage(multiplexer.PermissionDenied)
		return
	}
	if len(context.Fields) < 2 {
		context.SendMessage(multiplexer.InvalidArgument)
		return
	}
	switch context.Fields[1] {
	case "clear":
		if len(context.Fields) != 4 {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		index, err := strconv.Atoi(context.Fields[3])
		if err != nil {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		index = index - 1
		if index < 0 || index > 24 {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		member := context.GetMember(context.Fields[2])
		if member == nil {
			context.SendMessage(multiplexer.MissingUser)
			return
		}
		body, err := db.GetWarning(member.GuildID, member.User.ID)
		if !context.HandleError(err) {
			return
		}
		if body == "" {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		var warns []MemberWarning
		err = json.Unmarshal([]byte(body), &warns)
		if !context.HandleError(err) {
			return
		}
		if index > len(warns) {
			context.SendMessage(multiplexer.InvalidArgument)
			return
		}
		n := append(warns[:index], warns[index+1:]...)
		b, err := json.Marshal(n)
		if !context.HandleError(err) {
			return
		}
		err = db.SetWarning(member.GuildID, member.User.ID, string(b))
		if !context.HandleError(err) {
			return
		}
		context.SendMessage(fmt.Sprintf("Successfully cleared warning number %v.", index+1))
	default:
		member := context.GetMember(context.Fields[1])
		if member == nil {
			context.SendMessage(multiplexer.MissingUser)
			return
		}
		body, err := db.GetWarning(member.GuildID, member.User.ID)
		if !context.HandleError(err) {
			return
		}
		if body == "" {
			body = "[]"
		}
		var warns []MemberWarning
		err = json.Unmarshal([]byte(body), &warns)
		if !context.HandleError(err) {
			return
		}
		switch len(context.Fields) {
		case 2:
			embed := embedutil.New("Warnings", "List of warnings against "+member.User.Username)
			embed.Color = multiplexer.KappaColor
			for index, warn := range warns {
				embed.AddField(fmt.Sprintf("Warning on %s (%v)",
					warn.Time.UTC().Format("Mon Jan 2 15:04:05 2006"), index+1),
					fmt.Sprintf("Issuer: <@%s> \nReason: [%s](https://discord.com/channels/%s/%s/%s)",
						warn.UserID,
						warn.Text,
						warn.GuildID,
						warn.ChannelID,
						warn.MessageID), false)
			}
			context.SendEmbed("", embed)
		default:
			if len(warns) == 25 {
				context.SendMessage("Limit of 25 warnings per user reached, please clear some warnings and try again.")
				return
			}
			err = context.Session.ChannelMessageDelete(context.Channel.ID, context.Message.ID)
			if !context.HandleError(err) {
				return
			}
			message := context.StitchFields(2)
			m := context.SendMessage(fmt.Sprintf("Warning issued against %s with the reason `%s`.", member.Mention(), message))
			if m == nil {
				return
			}
			warns = append(warns, MemberWarning{
				Text:      message,
				Time:      time.Now(),
				GuildID:   context.Guild.ID,
				ChannelID: m.ChannelID,
				MessageID: m.ID,
				UserID:    context.User.ID,
			})
			b, err := json.Marshal(warns)
			if !context.HandleError(err) {
				return
			}
			err = db.SetWarning(member.GuildID, member.User.ID, string(b))
			if !context.HandleError(err) {
				return
			}
		}
	}
}

func ban(context *multiplexer.Context) {
	// Guild only
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
		return
	}

	// Has permission
	if !context.HasPermission(discordgo.PermissionBanMembers) {
		context.SendMessage(multiplexer.PermissionDenied)
		return
	}

	query := context.StitchFields(1)
	if query == "" {
		context.SendMessage(multiplexer.InvalidArgument)
		return
	}
	err = context.Ban(query)
	if err == discordgo.ErrUnauthorized {
		context.SendMessage(multiplexer.LackingPermission)
		return
	}
	if err == multiplexer.ErrUserNotFound {
		context.SendMessage(multiplexer.MissingUser)
		return
	}
	if !context.HandleError(err) {
		return
	}
	context.SendMessage("Successfully performed ban on specified user.")
}
