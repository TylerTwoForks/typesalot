package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/TylerTwoForks/typesalot/internal/db"
	"github.com/TylerTwoForks/typesalot/internal/webserver"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	conn, err := sql.Open("sqlite3", "./mydb.sqlite")
	if err != nil {
		log.Logger.Fatal()
	}

	q := db.New(conn)

	users, err := q.ListUsers(ctx())
	if err != nil {
		log.Logger.Fatal().Msg(err.Error())
	}

	log.Logger.Info().Interface("users", users).Msg("Fetched users")

	e := webserver.EchoServer()
	e.Static("/assets", "web/assets")

	go func() {
		if err := e.Start(":1323"); err != nil {
			log.Logger.Error().Msgf("Error starting server: %v", err)
		}
	}()

	gracefulShutdown(e)
}

func ctx() context.Context {
	return context.TODO() // or use context.WithTimeout()
}

func gracefulShutdown(e *echo.Echo) {
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // catch Ctrl+C and termination
	<-quit

	fmt.Println("Gracefully shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown: ", err)
	}
}
