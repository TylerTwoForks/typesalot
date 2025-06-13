package webserver

import (
	"github.com/TylerTwoForks/typesalot/web/templates"
	"github.com/labstack/echo/v4"
)

type EntryH struct {
}

func (eh *EntryH) NewHandler() *EntryH {
	return &EntryH{}
}

func (eh *EntryH) EntryRoutes(g *echo.Group) {
	eg := g.Group("entry")
	eg.GET("/test", func(c echo.Context) error {
		return Render(c, 200, templates.Test())
	})

}
