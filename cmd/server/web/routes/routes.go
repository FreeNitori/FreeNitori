package routes

import (
	"github.com/gin-gonic/gin"
)

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

type WebRoute struct {
	Pattern  string
	Handlers []gin.HandlerFunc
}
