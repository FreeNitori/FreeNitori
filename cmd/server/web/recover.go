package web

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			p := recover()
			if p != nil {
				log.Errorf("Panic occurred in web server, %s", p)
				context.JSON(http.StatusInternalServerError, datatypes.H{
					"error": p,
				})
			}
		}()
		context.Next()
	}
}
