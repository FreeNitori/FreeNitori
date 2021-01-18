package routes

import "gopkg.in/macaron.v1"

var (
	GetRoutes     []WebRoute
	PostRoutes    []WebRoute
	ComboRoutes   []WebRoute
	DeleteRoutes  []WebRoute
	HeadRoutes    []WebRoute
	OptionsRoutes []WebRoute
	PatchRoutes   []WebRoute
	PutRoutes     []WebRoute
	AnyRoutes     []WebRoute
)

type WebRoute struct {
	Pattern  string
	Handlers []macaron.Handler
}
