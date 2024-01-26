package clients

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

var SpotifyProvider goth.Provider

func TopArtists(c echo.Context) []byte {
	session, _ := session.Get("session", c)
	userSession := session.Values["user"].(goth.User)

	now := time.Now()
	if !userSession.ExpiresAt.After(now) {
		log.Println("need Refresh")
		t, err := SpotifyProvider.RefreshToken(userSession.RefreshToken)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
		}
		userSession.AccessToken = t.AccessToken
		userSession.ExpiresAt = t.Expiry
		session.Values["user"] = userSession
		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			log.Println("Error saving session", err)
		}

	}

	log.Println(now)
	log.Println(userSession.ExpiresAt)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/top/artists", nil)

	req.Header.Add("Authorization", "Bearer "+userSession.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error getting me")
	}
	defer resp.Body.Close()
	log.Println(resp.Status)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error")
	}
	return bodyBytes
}
