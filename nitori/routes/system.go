package routes

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"os"
	"runtime"
	"strconv"
	"time"
)

func init() {
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "about",
		AliasPatterns: []string{"info", "kappa", "information"},
		Description:   "Display system information.",
		Category:      multiplexer.SystemCategory,
		Handler:       about,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "stats",
		AliasPatterns: []string{},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler:       stats,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "invite",
		AliasPatterns: []string{"authorize", "oauth"},
		Description:   "Display authorization URL.",
		Category:      multiplexer.SystemCategory,
		Handler:       invite,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "reset-guild",
		AliasPatterns: []string{},
		Description:   "Reset current guild configuration.",
		Category:      multiplexer.SystemCategory,
		Handler:       resetGuild,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "shutdown",
		AliasPatterns: []string{"poweroff", "reboot", "restart"},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler:       shutdown,
	})
}

func about(context *multiplexer.Context) {
	embed := embedutil.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = vars.KappaColor
	embed.AddField("Homepage", config.Config.WebServer.BaseURL, true)
	embed.AddField("Version", state.Version(), true)
	embed.AddField("Commit Hash", state.Revision(), true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
	if vars.Administrator != nil {
		embed.AddField("Administrator", vars.Administrator.Username+"#"+vars.Administrator.Discriminator, true)
	}
	switch len(vars.Operator) {
	case 0:
	case 1:
		embed.AddField("Operator", vars.Operator[0].Username+"#"+vars.Operator[0].Discriminator, true)
	default:
		var usernames string
		for i, user := range vars.Operator {
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
	embed.SetFooter("FreeNitori Backend", "https://freenitori.jp/static/icon.min.png")
	context.SendEmbed(embed)
}

func stats(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(vars.OperatorOnly)
		return
	}

	var embed *embedutil.Embed

	uptime := state.Uptime()
	numGoroutine := runtime.NumGoroutine()

	embed = embedutil.NewEmbed("System Stats", "")
	embed.Color = vars.KappaColor
	embed.AddField("PID", strconv.Itoa(os.Getpid()), true)
	embed.AddField("Uptime", uptime.String(), true)
	embed.AddField("Goroutines", strconv.Itoa(numGoroutine), true)
	context.SendEmbed(embed)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memAllocated := sizeInMiB(memStats.Alloc)
	memTotal := sizeInMiB(memStats.TotalAlloc)
	memSys := sizeInMiB(memStats.Sys)
	lookups := memStats.Lookups
	memMallocs := memStats.Mallocs
	memFrees := memStats.Frees

	embed = embedutil.NewEmbed("", "")
	embed.Color = vars.KappaColor
	embed.AddField("Current Memory Allocated", memAllocated, true)
	embed.AddField("Total Memory Allocated", memTotal, true)
	embed.AddField("System Reported Allocation", memSys, true)
	embed.AddField("Pointer Lookups", strconv.Itoa(int(lookups)), true)
	embed.AddField("Memory Allocations", strconv.Itoa(int(memMallocs)), true)
	embed.AddField("Memory Frees", strconv.Itoa(int(memFrees)), true)
	context.SendEmbed(embed)

	heapAlloc := sizeInMiB(memStats.HeapAlloc)
	heapSys := sizeInMiB(memStats.HeapSys)
	heapIdle := sizeInMiB(memStats.HeapIdle)
	heapInuse := sizeInMiB(memStats.HeapInuse)
	heapReleased := sizeInMiB(memStats.HeapReleased)
	heapObjects := memStats.HeapObjects

	embed = embedutil.NewEmbed("", "")
	embed.Color = vars.KappaColor
	embed.AddField("Current Heap Usage", heapAlloc, true)
	embed.AddField("System Reported Usage", heapSys, true)
	embed.AddField("Heap Idle", heapIdle, true)
	embed.AddField("Heap In Use", heapInuse, true)
	embed.AddField("Heap Released", heapReleased, true)
	embed.AddField("Heap Objects", strconv.Itoa(int(heapObjects)), true)
	context.SendEmbed(embed)

	stackInuse := sizeInMiB(memStats.StackInuse)
	stackSys := sizeInMiB(memStats.StackSys)
	mspanInuse := sizeInMiB(memStats.MSpanInuse)
	mspanSys := sizeInMiB(memStats.MSpanSys)
	mcacheInuse := sizeInMiB(memStats.MCacheInuse)
	mcacheSys := sizeInMiB(memStats.MCacheSys)
	gcSys := sizeInMiB(memStats.GCSys)
	buckHashSys := sizeInMiB(memStats.BuckHashSys)
	otherSys := sizeInMiB(memStats.OtherSys)

	embed = embedutil.NewEmbed("", "")
	embed.Color = vars.KappaColor
	embed.AddField("Bootstrap Stack Usage", stackInuse, true)
	embed.AddField("System Reported Stack", stackSys, true)
	embed.AddField("MSpan Structures Usage", mspanInuse, true)
	embed.AddField("System Reported MSpan", mspanSys, true)
	embed.AddField("MCache Structures Usage", mcacheInuse, true)
	embed.AddField("System Reported MCache", mcacheSys, true)
	embed.AddField("GC Metadata Size", gcSys, true)
	embed.AddField("Profiling Bucket Hash Tables Size", buckHashSys, true)
	embed.AddField("Miscellaneous Off-heap Allocations", otherSys, true)
	context.SendEmbed(embed)

	nextGC := sizeInMiB(memStats.NextGC)
	lastGC := fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(memStats.LastGC))/1000/1000/1000)
	pauseTotalNs := fmt.Sprintf("%.1fs", float64(memStats.PauseTotalNs)/1000/1000/1000)
	pauseNs := fmt.Sprintf("%.3fs", float64(memStats.PauseNs[(memStats.NumGC+255)%256])/1000/1000/1000)
	numGC := memStats.NumGC

	embed = embedutil.NewEmbed("", "")
	embed.Color = vars.KappaColor
	embed.AddField("Next GC Recycle", nextGC, true)
	embed.AddField("Time Since Last GC", lastGC, true)
	embed.AddField("Total GC Pause", pauseTotalNs, true)
	embed.AddField("Last GC Pause", pauseNs, true)
	embed.AddField("Number of GCs", strconv.Itoa(int(numGC)), true)
	context.SendEmbed(embed)
}

func resetGuild(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(vars.OperatorOnly)
		return
	}
	config.ResetGuild(context.Guild.ID)
	context.SendMessage("Guild configuration has been reset.")
}

func shutdown(context *multiplexer.Context) {
	if !context.IsAdministrator() {
		context.SendMessage(vars.AdminOnly)
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
	embed := embedutil.NewEmbed("Invite", fmt.Sprintf("Click [this](%s) to invite Nitori.", state.InviteURL))
	context.SendEmbed(embed)
}

func sizeInMiB(size uint64) string {
	return fmt.Sprintf("%.2f MiB", float64(size)/1048576)
}
