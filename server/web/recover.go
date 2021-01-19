package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	"gopkg.in/macaron.v1"
	"net/http"
)

func recovery() macaron.Handler {
	return func(context *macaron.Context) {
		defer func() {
			p := recover()
			if p != nil {
				log.Errorf("Panic occurred in web server, %s", p)
				context.HTML(http.StatusInternalServerError, "error", datatypes.H{
					"Title":    datatypes.InternalServerError,
					"Subtitle": "PANIC!!! Something terribly wrong has occurred!!!",
					"Message":  p,
				})
			}
		}()
		context.Next()
	}
}
