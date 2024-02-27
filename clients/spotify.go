package clients

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anthonymq/go-stack-demo/logger"
	"github.com/anthonymq/go-stack-demo/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

var SpotifyProvider goth.Provider

func TopArtists(u goth.User) model.SpotifyTopArtists {
	resp, _ := callSpotifyApi("GET", "https://api.spotify.com/v1/me/top/artists", u)
	defer resp.Body.Close()
	results, _ := unmarshal[model.SpotifyTopArtists](resp)
	return results
}

func callSpotifyApi(method string, url string, u goth.User) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		logger.Get().Error(err.Error())
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+u.AccessToken)
	return client.Do(req)
}

func RenewAccessTokenIfExpired(c echo.Context, s sessions.Session, u goth.User) (goth.User, error) {
	if u.ExpiresAt.Before(time.Now()) {
		log.Println("need Refresh")
		t, err := SpotifyProvider.RefreshToken(u.RefreshToken)
		if err != nil {
			log.Println("Couldn't get a new AccessToken with this refreshToken")
			return goth.User{}, c.Redirect(http.StatusSeeOther, "/login")
		}
		u.AccessToken = t.AccessToken
		u.ExpiresAt = t.Expiry
		s.Values["user"] = u
		err = s.Save(c.Request(), c.Response().Writer)
		if err != nil {
			log.Println("Error saving session", err)
		}
	}
	return u, nil
}

func SearchTrack(u goth.User, query string) model.SpotifySearchResult {
	resp, err := callSpotifyApi("GET",
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track&limit=10", query),
		u,
	)
	defer resp.Body.Close()
	if err != nil {
		logger.Get().Error(err.Error())
	}
	logger.Get().Info(resp.Status)
	results, err := unmarshal[model.SpotifySearchResult](resp)
	if err != nil {
		logger.Get().Error(err.Error())
	}
	spew.Dump(results)
	return results
}
