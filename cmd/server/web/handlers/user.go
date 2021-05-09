package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/db"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/internals"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/sessioning"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func init() {
	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/api/user/:uid",
			Handlers: []gin.HandlerFunc{apiUser},
		},
		routes.WebRoute{
			Pattern:  "/api/guild/:gid",
			Handlers: []gin.HandlerFunc{apiGuild},
		},
		routes.WebRoute{
			Pattern:  "/api/guild/:gid/:key",
			Handlers: []gin.HandlerFunc{apiGuildKey},
		},
	)
}

func apiUser(context *gin.Context) {
	snowflake := context.Param("uid")
	if !config.ValidateSnowflake(snowflake) {
		context.JSON(http.StatusBadRequest, datatypes.H{
			"error": "invalid snowflake",
		})
		return
	}
	user, err := sessioning.FetchUser(snowflake)
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
		CreationTime:  time.Unix(int64(((func() (id int) { id, _ = strconv.Atoi(user.ID); return }()>>22)+1420070400000)/1000), 0).UTC(),
		Bot:           user.Bot,
	})
}

func apiGuild(context *gin.Context) {
	snowflake := context.Param("gid")
	if !config.ValidateSnowflake(snowflake) {
		context.JSON(http.StatusBadRequest, datatypes.H{
			"error": "invalid snowflake",
		})
		return
	}
	guild := sessioning.FetchGuild(snowflake)
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
		Name:         guild.Name,
		ID:           guild.ID,
		CreationTime: config.CreationTime(snowflake),
		IconURL:      guild.IconURL(),
		Members:      members,
	})
}

func apiGuildKey(context *gin.Context) {
	snowflake := context.Param("gid")
	if !config.ValidateSnowflake(snowflake) {
		context.JSON(http.StatusBadRequest, datatypes.H{
			"error": "invalid snowflake",
		})
		return
	}
	guild := sessioning.FetchGuild(snowflake)
	if guild == nil {
		context.JSON(http.StatusNotFound, datatypes.H{})
		return
	}
	switch context.Param("key") {
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
		expEnabled, err := db.ExpEnabled(guild.ID)
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
			expData, err := db.GetMemberExp(member.User, guild)
			if err != nil {
				context.JSON(http.StatusInternalServerError, datatypes.H{})
				return
			}
			levelData := internals.ExpToLevel(expData)
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
		context.JSON(http.StatusOK, datatypes.H{
			"Leaderboard": leaderboard,
			"GuildInfo": datatypes.GuildInfo{
				Name:    guild.Name,
				ID:      guild.ID,
				IconURL: guild.IconURL(),
				Members: nil,
			},
		})
	default:
		context.JSON(http.StatusNotFound, datatypes.H{
			"error": "not found",
		})
	}
}
