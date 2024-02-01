package clients

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

var SpotifyProvider goth.Provider

func TopArtists(u goth.User) []byte {
	// userSession := RenewAccessTokenIfExpired(c)
	return callSpotifyApi(u, "GET", "https://api.spotify.com/v1/me/top/artists")
}

func callSpotifyApi(u goth.User, method string, url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	req.Header.Add("Authorization", "Bearer "+u.AccessToken)
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

func RenewAccessTokenIfExpired(c echo.Context, s sessions.Session, u goth.User) goth.User {
	if u.ExpiresAt.Before(time.Now()) {
		log.Println("need Refresh")
		t, err := SpotifyProvider.RefreshToken(u.RefreshToken)
		if err != nil {
			log.Println("Couldn't get a new AccessToken with this refreshToken")
			c.Redirect(http.StatusSeeOther, "/login")
		}
		u.AccessToken = t.AccessToken
		u.ExpiresAt = t.Expiry
		s.Values["user"] = u
		err = s.Save(c.Request(), c.Response().Writer)
		if err != nil {
			log.Println("Error saving session", err)
		}
	}
	return u
}

// func (u goth.User) isExpired() bool {
// 	return u.ExpiresAt.Before(time.Now())
// }

func SearchTrack(u goth.User, query string) []byte {
	return callSpotifyApi(u, "GET", fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=5", query))
}
