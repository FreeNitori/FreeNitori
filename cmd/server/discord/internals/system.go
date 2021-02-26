package internals

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
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
		Handler:       stats,
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
	embed.AddField("Homepage", config.Config.WebServer.BaseURL, true)
	embed.AddField("Version", state.Version(), true)
	embed.AddField("Commit Hash", state.Revision(), true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
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

func stats(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}

	stats := config.Stats()

	var embed embedutil.Embed

	embed = embedutil.New("System Stats", "")
	embed.Color = multiplexer.KappaColor
	embed.AddField("PID", strconv.Itoa(stats.Process.PID), true)
	embed.AddField("Uptime", stats.Process.Uptime.Truncate(time.Second).String(), true)
	embed.AddField("Goroutines", strconv.Itoa(stats.Process.NumGoroutine), true)
	embed.AddField("Database Size", strconv.Itoa(int(stats.Process.DBSize)), true)

	embed.AddField("Intents", strconv.Itoa(stats.Discord.Intents), true)
	embed.AddField("Sharding", strconv.FormatBool(stats.Discord.Sharding), true)
	embed.AddField("Shards", strconv.Itoa(stats.Discord.Shards), true)
	embed.AddField("Guilds", strconv.Itoa(stats.Discord.Guilds), true)

	embed.AddField("Operating System", stats.Platform.GOOS, true)
	embed.AddField("Architecture", stats.Platform.GOARCH, true)
	embed.AddField("Go Root", stats.Platform.GOROOT, true)
	embed.AddField("Go Version", stats.Platform.GoVersion, true)
	context.SendEmbed("", embed)

	embed = embedutil.New("", "")
	embed.Color = multiplexer.KappaColor
	embed.AddField("Current Memory Allocated", stats.Mem.Allocated, true)
	embed.AddField("Total Memory Allocated", stats.Mem.Total, true)
	embed.AddField("System Reported Allocation", stats.Mem.Sys, true)
	embed.AddField("Pointer Lookups", strconv.Itoa(int(stats.Mem.Lookups)), true)
	embed.AddField("Memory Allocations", strconv.Itoa(int(stats.Mem.Mallocs)), true)
	embed.AddField("Memory Frees", strconv.Itoa(int(stats.Mem.Frees)), true)

	embed.AddField("Current Heap Usage", stats.Heap.Alloc, true)
	embed.AddField("System Reported Usage", stats.Heap.Sys, true)
	embed.AddField("Heap Idle", stats.Heap.Idle, true)
	embed.AddField("Heap In Use", stats.Heap.Inuse, true)
	embed.AddField("Heap Released", stats.Heap.Released, true)
	embed.AddField("Heap Objects", strconv.Itoa(int(stats.Heap.Objects)), true)

	embed.AddField("Bootstrap Stack Usage", stats.Misc.StackInuse, true)
	embed.AddField("System Reported Stack", stats.Misc.StackSys, true)
	embed.AddField("MSpan Structures Usage", stats.Misc.MSpanInuse, true)
	embed.AddField("System Reported MSpan", stats.Misc.MSpanSys, true)
	embed.AddField("MCache Structures Usage", stats.Misc.MCacheInuse, true)
	embed.AddField("System Reported MCache", stats.Misc.MCacheSys, true)
	embed.AddField("GC Metadata Size", stats.Misc.GCSys, true)
	embed.AddField("Profiling Bucket Hash Tables Size", stats.Misc.BuckHashSys, true)
	embed.AddField("Miscellaneous Off-heap Allocations", stats.Misc.OtherSys, true)
	context.SendEmbed("", embed)

	embed = embedutil.New("", "")
	embed.Color = multiplexer.KappaColor
	embed.AddField("Next GC Recycle", stats.GC.NextGC, true)
	embed.AddField("Time Since Last GC", stats.GC.LastGC, true)
	embed.AddField("Total GC Pause", stats.GC.PauseTotalNs, true)
	embed.AddField("Last GC Pause", stats.GC.PauseNs, true)
	embed.AddField("Number of GCs", strconv.Itoa(int(stats.GC.NumGC)), true)
	context.SendEmbed("", embed)
}

func resetGuild(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(multiplexer.OperatorOnly)
		return
	}
	config.ResetGuild(context.Guild.ID)
	context.SendMessage("Guild configuration has been reset.")
}

func shutdown(context *multiplexer.Context) {
	if !context.IsAdministrator() {
		context.SendMessage(multiplexer.AdminOnly)
		return
	}
	if map[string]bool{"reboot": true, "restart": true, "shutdown": false, "poweroff": false}[context.Fields[0]] {
		context.SendMessage("Attempting restart...")
		state.ExitCode <- -1
		return
	}
	context.SendMessage("Performing complete shutdown.")
	log.Info("Shutdown requested via Discord command.")
	state.ExitCode <- 0
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
	err = state.RawSession.UpdateGameStatus(0, config.Config.Discord.Presence)
	if err != nil {
		log.Warnf("Unable to update presence, %s", err)
	}

	log.Debugf("Session %s ready.",
		ready.SessionID)
}

func advanceCounter(context *multiplexer.Context) {
	if context.User.ID == context.Session.State.User.ID {
		return
	}
	err = config.AdvanceTotalMessages()
	if err != nil {
		log.Errorf("Unable to advance message counter, %s", err)
	}
}
