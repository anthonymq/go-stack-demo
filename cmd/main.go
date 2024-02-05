package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/handler"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

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
var SPOTIFY_CLIENT_KEY string
var SPOTIFY_CLIENT_SECRET string
var DEEZER_CLIENT_KEY string
var DEEZER_CLIENT_SECRET string

func main() {
	err := godotenv.Load("secrets.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	SPOTIFY_CLIENT_KEY = os.Getenv("SPOTIFY_CLIENT_KEY")
	SPOTIFY_CLIENT_SECRET = os.Getenv("SPOTIFY_CLIENT_SECRET")
	DEEZER_CLIENT_KEY = os.Getenv("DEEZER_CLIENT_KEY")
	DEEZER_CLIENT_SECRET = os.Getenv("DEEZER_CLIENT_SECRET")

	gothic.Store = &cookieStore

	clients.DeezerProvider = deezer.New(
		DEEZER_CLIENT_KEY,
		DEEZER_CLIENT_SECRET,
		"http://localhost:3000/auth/deezer/callback",
		"email",
	)
	clients.SpotifyProvider = spotify.New(
		SPOTIFY_CLIENT_KEY,
		SPOTIFY_CLIENT_SECRET,
		"http://localhost:3000/auth/spotify/callback",
		"user-top-read",
	)
	goth.UseProviders(clients.DeezerProvider)

	app := echo.New()
	// @TODO store auth in JWT Cookie
	app.Use(session.Middleware(&cookieStore))
	userHandler := handler.UserHandler{}
	loginHandler := handler.LoginHandler{}
	playlistHandler := handler.PlaylistHandler{}

	app.GET("/login", loginHandler.HandleLoginShow)
	app.GET("/auth/spotify/callback", func(c echo.Context) error {
		user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
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
		return c.Redirect(http.StatusSeeOther, "/app/user")
	})

	app.GET("/auth/deezer/callback", func(c echo.Context) error {
		code := c.QueryParam("code")
		authUrl := fmt.Sprintf("https://connect.deezer.com/oauth/access_token.php?app_id=%s&secret=%s&code=%s", DEEZER_CLIENT_KEY, DEEZER_CLIENT_SECRET, code)
		requestAccessTokenResponse, err := http.Get(authUrl)
		requestAccessTokenBody, err := io.ReadAll(requestAccessTokenResponse.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
		}
		defer requestAccessTokenResponse.Body.Close()
		vals, err := url.ParseQuery(string(requestAccessTokenBody))
		if err != nil {
			log.Println("no access token")
		}
		expires, err := strconv.Atoi(vals.Get("expires"))
		if err != nil {
			log.Panicf("could not get expires: %s\n", err)
		}

		deezerSession := clients.DeezerSession{
			AccessToken: vals.Get("access_token"),
			ExpiresAt:   time.Now().Add(time.Duration(expires)),
		}

		userConnected, err := clients.FetchUser(deezerSession)
		if err != nil {
			log.Panicln("Error retrieving connected user", err)
		}

		currentSession, err := session.Get("session", c)
		if err != nil {
			log.Println("Error retrieving session", err)
		}
		currentSession.Values["user"] = userConnected
		err = currentSession.Save(c.Request(), c.Response().Writer)
		if err != nil {
			log.Println("Error saving session", err)
		}
		return c.Redirect(http.StatusSeeOther, "/app/user")
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

	app.GET("/auth/deezer", func(c echo.Context) error {
		gothic.GetProviderName = func(req *http.Request) (string, error) { return "deezer", nil }
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
			if userConnected.Provider == "spotify" {
				clients.RenewAccessTokenIfExpired(c, *sess, userConnected)
			}
			return next(c)
		}
	})
	protectedGroup.GET("/user", userHandler.HandleUserShow)
	protectedGroup.GET("/playlist", playlistHandler.HandlePlaylistShow)
	protectedGroup.GET("/playlist/searchTracks", playlistHandler.HandlePlaylistSearchTracks)

	app.Logger.Fatal(app.Start(":3000"))
	fmt.Println("it works")
}
