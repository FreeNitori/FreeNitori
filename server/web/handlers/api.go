package handlers

import (
	"encoding/json"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/oauth"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func init() {
	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/api",
			Handlers: []gin.HandlerFunc{api},
		},
		routes.WebRoute{
			Pattern:  "/api/info",
			Handlers: []gin.HandlerFunc{apiInfo},
		},
		routes.WebRoute{
			Pattern:  "/api/stats",
			Handlers: []gin.HandlerFunc{apiStats},
		},
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
		routes.WebRoute{
			Pattern:  "/api/auth",
			Handlers: []gin.HandlerFunc{apiAuth},
		},
		routes.WebRoute{
			Pattern:  "/api/auth/user",
			Handlers: []gin.HandlerFunc{apiAuthUser},
		},
	)
}

func api(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"status": "OK!",
	})
}

func apiInfo(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"nitori_version":  state.Version(),
		"nitori_revision": state.Revision(),
		"invite_url":      state.InviteURL,
	})
}

func apiStats(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"total_messages":  config.GetTotalMessages(),
		"guilds_deployed": strconv.Itoa(len(vars.RawSession.State.Guilds)),
	})
}

func apiUser(context *gin.Context) {
	user, err := discord.FetchUser(context.Param("uid"))
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

func apiGuild(context *gin.Context) {
	guild := discord.FetchGuild(context.Param("gid"))
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

func apiGuildKey(context *gin.Context) {
	guild := discord.FetchGuild(context.Param("gid"))
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

func apiAuth(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"authorized": oauth.GetToken(context) != nil,
	})
}

func apiAuthUser(context *gin.Context) {
	token := oauth.GetToken(context)
	if token == nil {
		context.JSON(http.StatusOK, datatypes.H{
			"authorized": false,
			"user":       datatypes.UserInfo{},
		})
		return
	}
	client := oauth.Client(context, oauthConf)
	response, err := client.Get(discordgo.EndpointUser("@me"))
	if err != nil {
		panic(err)
	}
	if response.StatusCode == http.StatusUnauthorized {
		oauth.RemoveToken(context)
		context.JSON(http.StatusOK, datatypes.H{
			"authorized": false,
			"user":       datatypes.UserInfo{},
		})
		return
	}

	var user discordgo.User
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, datatypes.H{
		"authorized": true,
		"user": datatypes.UserInfo{
			Name:          user.Username,
			ID:            user.ID,
			AvatarURL:     user.AvatarURL("4096"),
			Discriminator: user.Discriminator,
			CreationTime:  time.Unix(int64(((func() (id int) { id, _ = strconv.Atoi(user.ID); return }()>>22)+1420070400000)/1000), 0),
			Bot:           user.Bot,
		},
	})
}
