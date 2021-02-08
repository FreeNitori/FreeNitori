package routes

import (
	"github.com/gin-gonic/gin"
)

// WebRoutes of specific requests.
var (
	GetRoutes     []WebRoute
	PostRoutes    []WebRoute
	DeleteRoutes  []WebRoute
	HeadRoutes    []WebRoute
	OptionsRoutes []WebRoute
	PatchRoutes   []WebRoute
	PutRoutes     []WebRoute
	AnyRoutes     []WebRoute
)

// WebRoute represents a route on the web server.
type WebRoute struct {
	Pattern  string
	Handlers []gin.HandlerFunc
}
