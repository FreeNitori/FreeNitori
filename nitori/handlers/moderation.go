package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func init() {
	ModerationCategory.Register(whois, "whois", []string{"lookup", "pfp"},
		"Lookup a user's detailed information by username, nickname or ID")
}

func whois(context *multiplexer.Context) {
	var user *discordgo.User
	var member *discordgo.Member

	// Just use the author if there's no arguments
	if len(context.Fields) == 1 {
		user = context.Author
		if !context.IsPrivate {
			member = context.Create.Member
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
		context.SendMessage(state.MissingUser)
		return
	}

	// Only generate the pfp stuff if that's what's required
	if context.Fields[0] == "pfp" {
		embed := formatter.NewEmbed("", "")
		embed.Color = context.Session.State.UserColor(user.ID, context.Create.ChannelID)
		embed.SetAuthor(user.Username, user.AvatarURL("128"))
		embed.SetImage(user.AvatarURL("4096"))
		context.SendEmbed(embed)
		return
	}

	// Make the message
	userID, err := strconv.Atoi(user.ID)
	if !context.HandleError(err) {
		return
	}
	creationTime := time.Unix(int64(((userID>>22)+1420070400000)/1000), 0)
	embed := formatter.NewEmbed("User Information", "")
	embed.Color = context.Session.State.UserColor(user.ID, context.Create.ChannelID)
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
	embed.AddField("Registration Date", creationTime.String(), true)
	if member != nil {
		joinTime, err := member.JoinedAt.Parse()
		if !context.HandleError(err) {
			return
		}
		embed.AddField("Join Date", joinTime.String(), true)
	}
	context.SendEmbed(embed)
}
