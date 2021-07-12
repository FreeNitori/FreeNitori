package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/db"
	discordRoutes "git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord/session"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord/snowflake"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/structs"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
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
	flake := context.Param("uid")
	if !snowflake.ValidateSnowflake(flake) {
		context.JSON(http.StatusBadRequest, structs.H{
			"error": "invalid snowflake",
		})
		return
	}
	user, err := session.FetchUser(flake)
	if err != nil {
		context.JSON(http.StatusInternalServerError, structs.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, structs.UserInfo{
		Name:          user.Username,
		ID:            user.ID,
		AvatarURL:     user.AvatarURL("4096"),
		Discriminator: user.Discriminator,
		CreationTime:  snowflake.CreationTime(user.ID),
		Bot:           user.Bot,
	})
}

func apiGuild(context *gin.Context) {
	flake := context.Param("gid")
	if !snowflake.ValidateSnowflake(flake) {
		context.JSON(http.StatusBadRequest, structs.H{
			"error": "invalid snowflake",
		})
		return
	}
	guild := session.FetchGuild(flake)
	if guild == nil {
		context.JSON(http.StatusNotFound, structs.H{
			"error": "not found",
		})
		return
	}
	var members []structs.UserInfo
	for _, member := range guild.Members {
		userInfo := structs.UserInfo{
			Name:          member.User.Username,
			ID:            member.User.ID,
			AvatarURL:     member.User.AvatarURL("128"),
			Discriminator: member.User.Discriminator,
			Bot:           member.User.Bot,
		}
		members = append(members, userInfo)
	}
	context.JSON(http.StatusOK, structs.GuildInfo{
		Name:         guild.Name,
		ID:           guild.ID,
		CreationTime: snowflake.CreationTime(flake),
		IconURL:      guild.IconURL(),
		Members:      members,
	})
}

func apiGuildKey(context *gin.Context) {
	flake := context.Param("gid")
	if !snowflake.ValidateSnowflake(flake) {
		context.JSON(http.StatusBadRequest, structs.H{
			"error": "invalid snowflake",
		})
		return
	}
	guild := session.FetchGuild(flake)
	if guild == nil {
		context.JSON(http.StatusNotFound, structs.H{})
		return
	}
	switch context.Param("key") {
	case "name":
		context.JSON(http.StatusOK, guild.Name)
	case "icon_url":
		context.JSON(http.StatusOK, guild.IconURL())
	case "members":
		var members []*structs.UserInfo
		for _, member := range guild.Members {
			userInfo := structs.UserInfo{
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
			context.JSON(http.StatusInternalServerError, structs.H{})
			return
		}
		if !expEnabled {
			context.JSON(http.StatusServiceUnavailable, structs.H{})
			return
		}
		var leaderboard []structs.LeaderboardEntry
		for _, member := range guild.Members {
			if member.User.Bot {
				continue
			}
			expData, err := db.GetMemberExp(member.User, guild)
			if err != nil {
				context.JSON(http.StatusInternalServerError, structs.H{})
				return
			}
			levelData := discordRoutes.ExpToLevel(expData)
			entry := structs.LeaderboardEntry{
				User: structs.UserInfo{
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
		context.JSON(http.StatusOK, structs.H{
			"Leaderboard": leaderboard,
			"GuildInfo": structs.GuildInfo{
				Name:    guild.Name,
				ID:      guild.ID,
				IconURL: guild.IconURL(),
				Members: nil,
			},
		})
	default:
		context.JSON(http.StatusNotFound, structs.H{
			"error": "not found",
		})
	}
}
