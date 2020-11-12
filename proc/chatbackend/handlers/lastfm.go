package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/state"
	"github.com/shkh/lastfm-go/lastfm"
	"regexp"
	"strconv"
)

func init() {
	AudioCategory.Register(fm, "fm", []string{"lastfm"}, "Query last song scrobbled to lastfm.")
	state.LastFM = lastfm.New(config.Config.LastFM.ApiKey, config.Config.LastFM.ApiSecret)
}

func fm(context *multiplexer.Context) {
	var username string
	switch len(context.Fields) {
	case 1:
	case 2:
		if context.Fields[1] == "unset" {
			err = config.ResetLastfm(context.Author, context.Guild)
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully reset lastfm username.")
			return
		}
		username = context.Fields[1]
	case 3:
		switch context.Fields[1] {
		case "set":
			if b, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, context.Fields[2]); !b || len(context.Fields[2]) < 2 || len(context.Fields[2]) > 15 {
				context.SendMessage(state.InvalidArgument)
				return
			}
			err = config.SetLastfm(context.Author, context.Guild, context.Fields[2])
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully set lastfm username to `" + context.Fields[2] + "`.")
			return
		case "unset":
			err = config.ResetLastfm(context.Author, context.Guild)
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully reset lastfm username.")
			return
		default:
			context.SendMessage(state.InvalidArgument)
			return
		}
	}
	if username == "" {
		username, err = config.GetLastfm(context.Author, context.Guild)
	}
	if !context.HandleError(err) {
		return
	}
	p := lastfm.P{"user": username, "limit": 1, "extended": 0}
	result, err := state.LastFM.User.GetRecentTracks(p)
	if err != nil {
		context.SendMessage("Please set your lastfm username `" + context.GenerateGuildPrefix() + "fm set <username>`.")
		return
	}
	if len(result.Tracks) < 1 {
		context.SendMessage("This username doesn't exist or does not have any scrobbles.")
		return
	}
	embed := embedutil.NewEmbed(result.Tracks[0].Name, result.Tracks[0].Artist.Name+" | "+result.Tracks[0].Album.Name)
	embed.SetAuthor(context.Author.Username, context.Author.AvatarURL("128"))
	embed.SetFooter(fmt.Sprintf("%s has %s scrobbles in total.", result.User, strconv.Itoa(result.Total)))
	embed.Color = context.Session.State.UserColor(context.Author.ID, context.Create.ChannelID)
	embed.URL = result.Tracks[0].Url
	if len(result.Tracks[0].Images) == 4 {
		embed.SetThumbnail(result.Tracks[0].Images[3].Url)
	}
	context.SendEmbed(embed)
}
