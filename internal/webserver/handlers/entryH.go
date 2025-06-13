package handlers

import "github.com/labstack/echo/v4"

type EntryH struct {
}

func (eh *EntryH) NewHandler() *EntryH {
	return &EntryH{}
}

func (eh *EntryH) EntryRoutes(g *echo.Group) {
	eg := g.Group("entry")
	eg.GET("/test", func(c echo.Context) error {
		return c.String(200, "essss")
	})

}
