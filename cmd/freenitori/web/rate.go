package web

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginDriver "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

var rateStore limiter.Store
var rateInstance *limiter.Limiter

func rateMiddleware() gin.HandlerFunc {
	// Create new memory-backed store
	rateStore = memory.NewStore()

	// Create new instance
	rateInstance = limiter.New(rateStore, limiter.Rate{
		Formatted: "",
		Period:    time.Duration(config.WebServer.RateLimitPeriod) * time.Second,
		Limit:     int64(config.WebServer.RateLimit),
	}, limiter.WithTrustForwardHeader(config.WebServer.ForwardedByClientIP))

	// Return middleware
	return ginDriver.NewMiddleware(rateInstance)
}
