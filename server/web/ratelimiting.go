package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

var (
	store       = memory.NewStore()
	rateLimiter = gin.NewMiddleware(instance)
	instance    = limiter.New(store, limiter.Rate{
		Formatted: "",
		Period:    time.Duration(config.Config.WebServer.RateLimitPeriod) * time.Second,
		Limit:     int64(config.Config.WebServer.RateLimit),
	})
)
