package web

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func RateLimiting(eventsLimit int, burstLimit int) gin.HandlerFunc {
	rateLimiter := rate.NewLimiter(rate.Limit(eventsLimit), burstLimit)
	return func(context *gin.Context) {
		if rateLimiter.Allow() {
			context.Next()
			return
		}
		context.AbortWithStatus(http.StatusTooManyRequests)
	}
}
