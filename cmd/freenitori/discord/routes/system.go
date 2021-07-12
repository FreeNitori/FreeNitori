package routes

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/db"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/stats"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func init() {
	state.Multiplexer.Ready = append(state.Multiplexer.Ready, setStatus)
	state.Multiplexer.MessageCreate = append(state.Multiplexer.MessageCreate, advanceCounter)
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "about",
		AliasPatterns: []string{"info", "kappa", "information"},
		Description:   "Display system information.",
		Category:      multiplexer.SystemCategory,
		Handler:       about,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "stats",
		AliasPatterns: []string{},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler:       stat,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "invite",
		AliasPatterns: []string{"authorize", "oauth"},
		Description:   "Display authorization URL.",
		Category:      multiplexer.SystemCategory,
		Handler:       invite,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "reset-guild",
		AliasPatterns: []string{},
		Description:   "Reset current guild configuration.",
		Category:      multiplexer.SystemCategory,
		Handler:       resetGuild,
	})
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "shutdown",
		AliasPatterns: []string{"poweroff", "reboot", "restart"},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler:       shutdown,
	})
}

func about(context *multiplexer.Context) {
	embed := embedutil.New(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = multiplexer.KappaColor
	embed.AddField("Homepage", config.WebServer.BaseURL, true)
	embed.AddField("Version", state.Version(), true)
	embed.AddField("Commit Hash", state.Revision(), true)
	embed.AddField("Processed Messages", strconv.Itoa(db.GetTotalMessages()), true)
	if state.Multiplexer.Administrator != nil {
		embed.AddField("Administrator", state.Multiplexer.Administrator.Username+"#"+state.Multiplexer.Administrator.Discriminator, true)
	}
	switch len(state.Multiplexer.Operator) {
	case 0:
	case 1:
		embed.AddField("Operator", state.Multiplexer.Operator[0].Username+"#"+state.Multiplexer.Operator[0].Discriminator, true)
	default:
		var usernames string
		for i, user := range state.Multiplexer.Operator {
			switch i {
			case 0:
				usernames += user.Username + "#" + user.Discriminator
			default:
				usernames += ", " + user.Username + "#" + user.Discriminator
			}
		}
		embed.AddField("Operators", usernames, true)
	}
	embed.SetThumbnail(context.Session.State.User.AvatarURL("256"))
	embed.SetFooter("FreeNitori Backend", "https://freenitori.jp/img/icon.min.png")
	context.SendEmbed("", embed)
}

func stat(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}

	info := stats.Get()

	var embed embedutil.Embed

	embed = embedutil.New("System Stats", "")
	embed.Color = multiplexer.KappaColor
	embed.AddField("PID", strconv.Itoa(info.Process.PID), true)
	embed.AddField("Uptime", info.Process.Uptime.Truncate(time.Second).String(), true)
	embed.AddField("Goroutines", strconv.Itoa(info.Process.NumGoroutine), true)
	embed.AddField("Database Size", strconv.Itoa(int(info.Process.DBSize)), true)

	embed.AddField("Intents", strconv.Itoa(info.Discord.Intents), true)
	embed.AddField("Sharding", strconv.FormatBool(info.Discord.Sharding), true)
	embed.AddField("Shards", strconv.Itoa(info.Discord.Shards), true)
	embed.AddField("Guilds", strconv.Itoa(info.Discord.Guilds), true)

	embed.AddField("Operating System", info.Platform.GOOS, true)
	embed.AddField("Architecture", info.Platform.GOARCH, true)
	embed.AddField("Go Root", info.Platform.GOROOT, true)
	embed.AddField("Go Version", info.Platform.GoVersion, true)
	context.SendEmbed("", embed)

	embed = embedutil.New("", "")
	embed.Color = multiplexer.KappaColor
	embed.AddField("Current Memory Allocated", info.Mem.Allocated, true)
	embed.AddField("Total Memory Allocated", info.Mem.Total, true)
	embed.AddField("System Reported Allocation", info.Mem.Sys, true)
	embed.AddField("Pointer Lookups", strconv.Itoa(int(info.Mem.Lookups)), true)
	embed.AddField("Memory Allocations", strconv.Itoa(int(info.Mem.Mallocs)), true)
	embed.AddField("Memory Frees", strconv.Itoa(int(info.Mem.Frees)), true)

	embed.AddField("Current Heap Usage", info.Heap.Alloc, true)
	embed.AddField("System Reported Usage", info.Heap.Sys, true)
	embed.AddField("Heap Idle", info.Heap.Idle, true)
	embed.AddField("Heap In Use", info.Heap.Inuse, true)
	embed.AddField("Heap Released", info.Heap.Released, true)
	embed.AddField("Heap Objects", strconv.Itoa(int(info.Heap.Objects)), true)

	embed.AddField("Bootstrap Stack Usage", info.Misc.StackInuse, true)
	embed.AddField("System Reported Stack", info.Misc.StackSys, true)
	embed.AddField("MSpan Structures Usage", info.Misc.MSpanInuse, true)
	embed.AddField("System Reported MSpan", info.Misc.MSpanSys, true)
	embed.AddField("MCache Structures Usage", info.Misc.MCacheInuse, true)
	embed.AddField("System Reported MCache", info.Misc.MCacheSys, true)
	embed.AddField("GC Metadata Size", info.Misc.GCSys, true)
	embed.AddField("Profiling Bucket Hash Tables Size", info.Misc.BuckHashSys, true)
	embed.AddField("Miscellaneous Off-heap Allocations", info.Misc.OtherSys, true)
	context.SendEmbed("", embed)

	embed = embedutil.New("", "")
	embed.Color = multiplexer.KappaColor
	embed.AddField("Next GC Recycle", info.GC.NextGC, true)
	embed.AddField("Time Since Last GC", info.GC.LastGC, true)
	embed.AddField("Total GC Pause", info.GC.PauseTotalNs, true)
	embed.AddField("Last GC Pause", info.GC.PauseNs, true)
	embed.AddField("Number of GCs", strconv.Itoa(int(info.GC.NumGC)), true)
	context.SendEmbed("", embed)
}

func resetGuild(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}
	db.ResetGuild(context.Guild.ID)
	context.SendMessage("Guild configuration has been reset.")
}

func shutdown(context *multiplexer.Context) {
	if !context.IsAdministrator() {
		context.SendMessage(multiplexer.AdminOnly)
		return
	}
	if map[string]bool{"reboot": true, "restart": true, "shutdown": false, "poweroff": false}[context.Fields[0]] {
		message := context.SendMessage("Attempting restart...")
		if message != nil {
			state.Reincarnation = message.ChannelID + "\t" + message.ID + "\t" + "Restart complete."
		}
		state.Exit <- -1
		return
	}
	context.SendMessage("Performing complete shutdown.")
	log.Info("Shutdown requested via Discord command.")
	state.Exit <- 0
}

func invite(context *multiplexer.Context) {
	embed := embedutil.New("Invite", fmt.Sprintf("Click [this](%s) to invite Nitori.", state.InviteURL))
	context.SendEmbed("", embed)
}

func setStatus(context *multiplexer.Context) {
	ready, ok := context.Event.(*discordgo.Ready)
	if !ok {
		return
	}
	err := state.Session.UpdateGameStatus(0, config.Discord.Presence)
	if err != nil {
		log.Warnf("Error updating presence, %s", err)
	}

	log.Debugf("Session %s ready.",
		ready.SessionID)
}

func advanceCounter(context *multiplexer.Context) {
	if context.User.ID == context.Session.State.User.ID {
		return
	}
	err := db.AdvanceTotalMessages()
	if err != nil {
		log.Errorf("Error advancing message counter, %s", err)
	}
}
