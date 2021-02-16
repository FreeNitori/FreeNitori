package config

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"os"
	"runtime"
	"time"
)

// SystemStats represents stats of the instance.
type SystemStats struct {
	Process struct {
		PID          int
		Uptime       time.Duration
		NumGoroutine int
		DBSize       int64
	}
	Platform struct {
		GoVersion string
		GOOS      string
		GOARCH    string
		GOROOT    string
	}
	Discord struct {
		Intents  int
		Sharding bool
		Shards   int
		Guilds   int
	}
	Mem struct {
		Allocated string
		Total     string
		Sys       string
		Lookups   uint64
		Mallocs   uint64
		Frees     uint64
	}
	Heap struct {
		Alloc    string
		Sys      string
		Idle     string
		Inuse    string
		Released string
		Objects  uint64
	}
	GC struct {
		NextGC       string
		LastGC       string
		PauseTotalNs string
		PauseNs      string
		NumGC        uint32
	}
	Misc struct {
		StackInuse  string
		StackSys    string
		MSpanInuse  string
		MSpanSys    string
		MCacheInuse string
		MCacheSys   string
		GCSys       string
		BuckHashSys string
		OtherSys    string
	}
}

// Stats returns a populated SystemStats.
func Stats() SystemStats {
	var systemStats SystemStats

	systemStats.Process.PID = os.Getpid()
	systemStats.Process.Uptime = state.Uptime()
	systemStats.Process.NumGoroutine = runtime.NumGoroutine()
	systemStats.Process.DBSize = database.Database.Size()

	systemStats.Platform.GoVersion = runtime.Version()
	systemStats.Platform.GOOS = runtime.GOOS
	systemStats.Platform.GOARCH = runtime.GOARCH
	systemStats.Platform.GOROOT = runtime.GOROOT()

	systemStats.Discord.Intents = int(state.RawSession.Identify.Intents)
	systemStats.Discord.Sharding = Config.Discord.Shard
	systemStats.Discord.Shards = len(state.ShardSessions)
	systemStats.Discord.Guilds = len(state.RawSession.State.Guilds)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	systemStats.Mem.Allocated = sizeInMiB(memStats.Alloc)
	systemStats.Mem.Total = sizeInMiB(memStats.TotalAlloc)
	systemStats.Mem.Sys = sizeInMiB(memStats.Sys)
	systemStats.Mem.Lookups = memStats.Lookups
	systemStats.Mem.Mallocs = memStats.Mallocs
	systemStats.Mem.Frees = memStats.Frees

	systemStats.Heap.Alloc = sizeInMiB(memStats.HeapAlloc)
	systemStats.Heap.Sys = sizeInMiB(memStats.HeapSys)
	systemStats.Heap.Idle = sizeInMiB(memStats.HeapIdle)
	systemStats.Heap.Inuse = sizeInMiB(memStats.HeapInuse)
	systemStats.Heap.Released = sizeInMiB(memStats.HeapReleased)
	systemStats.Heap.Objects = memStats.HeapObjects

	systemStats.Misc.StackInuse = sizeInMiB(memStats.StackInuse)
	systemStats.Misc.StackSys = sizeInMiB(memStats.StackSys)
	systemStats.Misc.MSpanInuse = sizeInMiB(memStats.MSpanInuse)
	systemStats.Misc.MSpanSys = sizeInMiB(memStats.MSpanSys)
	systemStats.Misc.MCacheInuse = sizeInMiB(memStats.MCacheInuse)
	systemStats.Misc.MCacheSys = sizeInMiB(memStats.MCacheSys)
	systemStats.Misc.GCSys = sizeInMiB(memStats.GCSys)
	systemStats.Misc.BuckHashSys = sizeInMiB(memStats.BuckHashSys)
	systemStats.Misc.OtherSys = sizeInMiB(memStats.OtherSys)

	systemStats.GC.NextGC = sizeInMiB(memStats.NextGC)
	systemStats.GC.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(memStats.LastGC))/1000/1000/1000)
	systemStats.GC.PauseTotalNs = fmt.Sprintf("%.1fs", float64(memStats.PauseTotalNs)/1000/1000/1000)
	systemStats.GC.PauseNs = fmt.Sprintf("%.3fs", float64(memStats.PauseNs[(memStats.NumGC+255)%256])/1000/1000/1000)
	systemStats.GC.NumGC = memStats.NumGC

	return systemStats
}

func sizeInMiB(size uint64) string {
	return fmt.Sprintf("%.2f MiB", float64(size)/1048576)
}
