package webserver

import (
	"fmt"
	"os"
	"time"

	"github.com/TylerTwoForks/typesalot/web/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// EchoServer(conn *gorm.DB)
func EchoServer() *echo.Echo {
	e := echo.New()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	e.Use(ZerologMiddleware(logger))

	// echoRoutes(e, conn)
	echoRoutes(e)

	//this is just logging the routes that we're using on server startup
	for _, route := range e.Routes() {
		log.Debug().Msg(fmt.Sprintf("Method: %s, Path: %s", route.Method, route.Path))
	}

	//base route. this sholud return whatever template we decide for the home page. Possibly a login page.
	e.GET("/", func(c echo.Context) error {
		return Render(c, 200, templates.Playground())
	})

	return e
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

// echoRoutes(e *echo.Echo, c *gorm.DB)
func echoRoutes(e *echo.Echo) {
	//grouping routes
	rg := e.Group("/") //route group (rg) - this is the default group.
	eh := EntryH{}
	eh.NewHandler().EntryRoutes(rg)
	/*
		 	lr := &repos.LaborRepo{DB: c}
			lh := LaborHandler{R: lr.NewRepo()}
			lh.NewHandler().LaborRoutes(rg)

			jr := &repos.JobRepo{DB: c}
			jh := JobHandler{R: jr.NewRepo()}
			jh.NewHandler().JobRoutesE(rg)
	*/
}

type Handler interface {
	NewHandler()
}

func ZerologMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			event := logger.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Dur("latency", stop.Sub(start)).
				Str("remote_ip", c.RealIP())

			if err != nil {
				c.Error(err)
				event.Err(err)
			}

			event.Msg("request handled")

			return err
		}
	}
}
