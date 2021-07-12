package web

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

var (
	rateStore      = memory.NewStore()
	rateMiddleware = gin.NewMiddleware(instance)
	instance       = limiter.New(rateStore, limiter.Rate{
		Formatted: "",
		Period:    time.Duration(config.WebServer.RateLimitPeriod) * time.Second,
		Limit:     int64(config.WebServer.RateLimit),
	}, limiter.WithTrustForwardHeader(config.WebServer.ForwardedByClientIP))
)
