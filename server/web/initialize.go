package web

import (
	"errors"
	"git.randomchars.net/RandomChars/FreeNitori/binaries/static"
	"git.randomchars.net/RandomChars/FreeNitori/binaries/tmpl"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/jsontypes"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Error messages and stuff
const (
	internalServerError   = "Internal Server Error"
	noSuchFileOrDirectory = "No such file or directory"
	badRequest            = "Bad Request"
	serviceUnavailable    = "Service Unavailable"
)

type leaderboardEntry struct {
	User       *jsontypes.UserInfo
	Experience int
	Level      int
}

func Initialize() error {
	// Initialize the engine
	gin.SetMode(gin.ReleaseMode)
	Engine = gin.New()
	Engine.Use(rateLimiter)
	Engine.ForwardedByClientIP = config.Config.WebServer.ForwardedByClientIP

	// Register templates
	templates := template.New("/")
	for _, path := range tmpl.AssetNames() {
		if strings.HasPrefix(path, "") {
			templateBin, _ := tmpl.Asset(path)
			templates, err = templates.New(path).Parse(string(templateBin))
			if err != nil {
				return errors.New("failed to parse templates")
			}
		}
	}
	Engine.SetHTMLTemplate(templates)

	// Register static files
	Engine.StaticFS("/static", static.AssetFile())
	Engine.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "error.html", gin.H{
			"Title":    noSuchFileOrDirectory,
			"Subtitle": "This route doesn't seem to exist.",
			"Message":  "I wonder how you got here...",
		})
	})

	// Register page routes
	Engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	Engine.GET("/lookup", func(context *gin.Context) {
		context.HTML(http.StatusOK, "lookup.html", nil)
	})
	Engine.GET("/guild/:gid/leaderboard", func(context *gin.Context) {
		guild := discord.FetchGuild(context.Param("gid"))
		if guild == nil {
			context.HTML(http.StatusNotFound, "error.html", gin.H{
				"Title":    noSuchFileOrDirectory,
				"Subtitle": "This guild doesn't seem to exist.",
				"Message":  "Maybe you got the wrong URL?",
			})
			return
		}
		expEnabled, err := config.ExpEnabled(guild.ID)
		if err != nil {
			context.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Title":    internalServerError,
				"Subtitle": "Failed to fetch experience system enablement status.",
				"Message":  "Nitori taking a nap?",
			})
			return
		}
		if !expEnabled {
			context.HTML(http.StatusServiceUnavailable, "error.html", gin.H{
				"Title":    serviceUnavailable,
				"Subtitle": "This feature is disabled in your guild.",
				"Message":  "Moderators don't like Nitori?",
			})
			return
		}
		context.HTML(http.StatusOK, "leaderboard.html", gin.H{
			"GuildName": guild.Name,
			"GuildIcon": guild.IconURL,
		})
	})

	// Register JSON API routes
	Engine.GET("/api", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "OK!"})
	})
	Engine.GET("/api/info", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"nitori_version":  state.Version(),
			"nitori_revision": state.Revision(),
			"invite_url":      state.InviteURL,
		})
	})
	Engine.GET("/api/stats", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"total_messages":  config.GetTotalMessages(),
			"guilds_deployed": strconv.Itoa(len(vars.RawSession.State.Guilds)),
		})
	})
	Engine.GET("/api/user/:uid", func(context *gin.Context) {
		user, err := discord.FetchUser(context.Param("uid"))
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, jsontypes.UserInfo{
			Name:          user.Username,
			ID:            user.ID,
			AvatarURL:     user.AvatarURL("4096"),
			Discriminator: user.Discriminator,
			CreationTime:  time.Unix(int64(((func() (id int) { id, _ = strconv.Atoi(user.ID); return }()>>22)+1420070400000)/1000), 0),
			Bot:           user.Bot,
		})
	})
	Engine.GET("/api/guild/:gid", func(context *gin.Context) {
		guild := discord.FetchGuild(context.Param("gid"))
		if guild == nil {
			context.JSON(http.StatusNotFound, gin.H{})
			return
		}
		var members []*jsontypes.UserInfo
		for _, member := range guild.Members {
			userInfo := jsontypes.UserInfo{
				Name:          member.User.Username,
				ID:            member.User.ID,
				AvatarURL:     member.User.AvatarURL("128"),
				Discriminator: member.User.Discriminator,
				Bot:           member.User.Bot,
			}
			members = append(members, &userInfo)
		}
		context.JSON(http.StatusOK, jsontypes.GuildInfo{
			Name:    guild.Name,
			ID:      guild.ID,
			IconURL: guild.IconURL(),
			Members: members,
		})
	})
	Engine.GET("/api/guild/:gid/:key", func(context *gin.Context) {
		guild := discord.FetchGuild(context.Param("gid"))
		if guild == nil {
			context.JSON(http.StatusNotFound, gin.H{})
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
			var members []*jsontypes.UserInfo
			for _, member := range guild.Members {
				userInfo := jsontypes.UserInfo{
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
				context.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			if !expEnabled {
				context.JSON(http.StatusServiceUnavailable, gin.H{})
				return
			}
			var leaderboard []*leaderboardEntry
			for _, member := range guild.Members {
				if member.User.Bot {
					continue
				}
				expData, err := config.GetMemberExp(member.User, guild)
				if err != nil {
					context.JSON(http.StatusInternalServerError, gin.H{})
					return
				}
				levelData := config.ExpToLevel(expData)
				entry := leaderboardEntry{
					User: &jsontypes.UserInfo{
						Name:          member.User.Username,
						ID:            member.User.ID,
						AvatarURL:     member.User.AvatarURL("128"),
						Discriminator: member.User.Discriminator,
						Bot:           member.User.Bot,
					},
					Experience: expData,
					Level:      levelData,
				}
				leaderboard = append(leaderboard, &entry)
			}
			sort.Slice(leaderboard, func(i, j int) bool {
				return leaderboard[i].Experience > leaderboard[j].Experience
			})
			context.JSON(http.StatusOK, leaderboard)
		default:
			context.JSON(http.StatusNotFound, gin.H{})
		}
	})
	return nil
}
