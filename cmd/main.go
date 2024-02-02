package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/handler"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

	// "github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/spotify"
)

var sessionKey = "mysecretsessionkey"
var cookieStore = *sessions.NewCookieStore([]byte(sessionKey))
var SessionName = "auth"
var SPOTIFY_CLIENT_KEY string
var SPOTIFY_CLIENT_SECRET string

func main() {
	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SPOTIFY_CLIENT_KEY = os.Getenv("SPOTIFY_CLIENT_KEY")
	SPOTIFY_CLIENT_SECRET = os.Getenv("SPOTIFY_CLIENT_SECRET")
	gothic.Store = &cookieStore
	clients.SpotifyProvider = spotify.New(
		SPOTIFY_CLIENT_KEY,
		SPOTIFY_CLIENT_SECRET,
		"http://localhost:3000/auth/spotify/callback",
		"user-top-read",
	)
	goth.UseProviders(clients.SpotifyProvider)
	app := echo.New()
	// @TODO store auth in JWT Cookie
	app.Use(session.Middleware(&cookieStore))
	userHandler := handler.UserHandler{}
	loginHandler := handler.LoginHandler{}
	playlistHandler := handler.PlaylistHandler{}
	app.GET("/login", loginHandler.HandleLoginShow)
	app.GET("/auth/spotify/callback", func(c echo.Context) error {

		user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
		// if err != nil {
		session, err := session.Get("session", c)
		if err != nil {
			log.Println("Error retrieving session", err)
		}
		user.RawData = map[string]interface{}{}
		session.Values["user"] = user
		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			log.Println("Error saving session", err)
		}
		println(spew.Sdump(session.Values["user"]))
		return c.Redirect(http.StatusSeeOther, "/app/user")
		// }
	})
	app.GET("/auth/spotify", func(c echo.Context) error {
		gothic.GetProviderName = func(req *http.Request) (string, error) { return "spotify", nil }
		if _, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
			log.Println("Already loggedInd")
			userHandler.HandleUserShow(c)
		} else {
			app.Logger.Error("ERR", err)
			log.Println("Not already loggedInd")
			gothic.BeginAuthHandler(c.Response(), c.Request())
		}
		return nil
	})
	protectedGroup := app.Group("/app", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			userConnected, exists := sess.Values["user"].(goth.User)
			if !exists {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "please provide valid credentials")
			}
			clients.RenewAccessTokenIfExpired(c, *sess, userConnected)
			return next(c)
		}
	})
	protectedGroup.GET("/user", userHandler.HandleUserShow)
	protectedGroup.GET("/playlist", playlistHandler.HandlePlaylistShow)
	protectedGroup.GET("/playlist/searchTracks", playlistHandler.HandlePlaylistSearchTracks)

	app.Logger.Fatal(app.Start(":3000"))
	fmt.Println("it works")

}
