package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"gopkg.in/macaron.v1"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func init() {
	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/api",
			Handlers: []macaron.Handler{api},
		},
		routes.WebRoute{
			Pattern:  "/api/info",
			Handlers: []macaron.Handler{apiInfo},
		},
		routes.WebRoute{
			Pattern:  "/api/stats",
			Handlers: []macaron.Handler{apiStats},
		},
		routes.WebRoute{
			Pattern:  "/api/user/:uid",
			Handlers: []macaron.Handler{apiUser},
		},
		routes.WebRoute{
			Pattern:  "/api/guild/:gid",
			Handlers: []macaron.Handler{apiGuild},
		},
		routes.WebRoute{
			Pattern:  "/api/guild/:gid/:key",
			Handlers: []macaron.Handler{apiGuildKey},
		},
	)
}

func api(context *macaron.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"status": "OK!",
	})
}

func apiInfo(context *macaron.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"nitori_version":  state.Version(),
		"nitori_revision": state.Revision(),
		"invite_url":      state.InviteURL,
	})
}

func apiStats(context *macaron.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"total_messages":  config.GetTotalMessages(),
		"guilds_deployed": strconv.Itoa(len(vars.RawSession.State.Guilds)),
	})
}

func apiUser(context *macaron.Context) {
	user, err := discord.FetchUser(context.Params("uid"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, datatypes.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, datatypes.UserInfo{
		Name:          user.Username,
		ID:            user.ID,
		AvatarURL:     user.AvatarURL("4096"),
		Discriminator: user.Discriminator,
		CreationTime:  time.Unix(int64(((func() (id int) { id, _ = strconv.Atoi(user.ID); return }()>>22)+1420070400000)/1000), 0),
		Bot:           user.Bot,
	})
}

func apiGuild(context *macaron.Context) {
	guild := discord.FetchGuild(context.Params("gid"))
	if guild == nil {
		context.JSON(http.StatusNotFound, datatypes.H{
			"error": "not found",
		})
		return
	}
	var members []datatypes.UserInfo
	for _, member := range guild.Members {
		userInfo := datatypes.UserInfo{
			Name:          member.User.Username,
			ID:            member.User.ID,
			AvatarURL:     member.User.AvatarURL("128"),
			Discriminator: member.User.Discriminator,
			Bot:           member.User.Bot,
		}
		members = append(members, userInfo)
	}
	context.JSON(http.StatusOK, datatypes.GuildInfo{
		Name:    guild.Name,
		ID:      guild.ID,
		IconURL: guild.IconURL(),
		Members: members,
	})
}

func apiGuildKey(context *macaron.Context) {
	guild := discord.FetchGuild(context.Params("gid"))
	if guild == nil {
		context.JSON(http.StatusNotFound, datatypes.H{})
		return
	}
	switch context.Params("key") {
	case "id":
		context.JSON(http.StatusOK, guild.ID)
	case "name":
		context.JSON(http.StatusOK, guild.Name)
	case "icon_url":
		context.JSON(http.StatusOK, guild.IconURL())
	case "members":
		var members []*datatypes.UserInfo
		for _, member := range guild.Members {
			userInfo := datatypes.UserInfo{
				Name:          member.User.Username,
				ID:            member.User.ID,
				AvatarURL:     member.User.AvatarURL("128"),
				Discriminator: member.User.Discriminator,
				Bot:           member.User.Bot,
			}
			members = append(members, &userInfo)
		}
		context.JSON(http.StatusOK, members)
	case "leaderboard":
		expEnabled, err := config.ExpEnabled(guild.ID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, datatypes.H{})
			return
		}
		if !expEnabled {
			context.JSON(http.StatusServiceUnavailable, datatypes.H{})
			return
		}
		var leaderboard []datatypes.LeaderboardEntry
		for _, member := range guild.Members {
			if member.User.Bot {
				continue
			}
			expData, err := config.GetMemberExp(member.User, guild)
			if err != nil {
				context.JSON(http.StatusInternalServerError, datatypes.H{})
				return
			}
			levelData := config.ExpToLevel(expData)
			entry := datatypes.LeaderboardEntry{
				User: datatypes.UserInfo{
					Name:          member.User.Username,
					ID:            member.User.ID,
					AvatarURL:     member.User.AvatarURL("128"),
					Discriminator: member.User.Discriminator,
					Bot:           member.User.Bot,
				},
				Experience: expData,
				Level:      levelData,
			}
			leaderboard = append(leaderboard, entry)
		}
		sort.Slice(leaderboard, func(i, j int) bool {
			return leaderboard[i].Experience > leaderboard[j].Experience
		})
		context.JSON(http.StatusOK, leaderboard)
	default:
		context.JSON(http.StatusNotFound, datatypes.H{
			"error": "not found",
		})
	}
}
