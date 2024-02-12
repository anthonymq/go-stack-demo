package main

import (
	"os"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/common"
	"github.com/anthonymq/go-stack-demo/handlers"
	"github.com/anthonymq/go-stack-demo/logger"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	// "github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/spotify"
)

var sessionKey = "mysecretsessionkey"
var cookieStore = *sessions.NewCookieStore([]byte(sessionKey))
var SessionName = "auth"

func main() {
	log := logger.Get()

	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	common.SPOTIFY_CLIENT_KEY = os.Getenv("SPOTIFY_CLIENT_KEY")
	common.SPOTIFY_CLIENT_SECRET = os.Getenv("SPOTIFY_CLIENT_SECRET")
	common.DEEZER_CLIENT_KEY = os.Getenv("DEEZER_CLIENT_KEY")
	common.DEEZER_CLIENT_SECRET = os.Getenv("DEEZER_CLIENT_SECRET")
	clients.DeezerProvider = deezer.New(
		common.DEEZER_CLIENT_KEY,
		common.DEEZER_CLIENT_SECRET,
		"http://localhost:3000/auth/deezer/callback",
		"email",
	)
	clients.SpotifyProvider = spotify.New(
		common.SPOTIFY_CLIENT_KEY,
		common.SPOTIFY_CLIENT_SECRET,
		"http://localhost:3000/auth/spotify/callback",
		"user-top-read",
	)

	gothic.Store = &cookieStore
	goth.UseProviders(clients.DeezerProvider, clients.SpotifyProvider)

	app := echo.New()
	// @TODO store auth in JWT Cookie
	app.Use(session.Middleware(&cookieStore))

	userHandler := handlers.UserHandler{}
	authHandler := handlers.AuthHandler{}
	playlistHandler := handlers.PlaylistHandler{}

	handlers.SetupRoutes(app,
		&authHandler,
		&userHandler,
		&playlistHandler,
	)

	err = app.Start(":3000")
	log.Fatal("Error starting server", zap.Error(err))
}
