package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	"github.com/gin-gonic/gin"
	"net/http"
)

func recovery() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			p := recover()
			if p != nil {
				log.Errorf("Panic occurred in web server, %s", p)
				context.HTML(http.StatusInternalServerError, "error.tmpl", datatypes.H{
					"Title":    datatypes.InternalServerError,
					"Subtitle": "PANIC!!! Something terribly wrong has occurred!!!",
					"Message":  p,
				})
			}
		}()
		context.Next()
	}
}
