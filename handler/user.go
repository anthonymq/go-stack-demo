package handler

import (
	"encoding/json"
	"log"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/model"
	"github.com/anthonymq/go-stack-demo/view/user"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

type UserHandler struct {
}

func (h UserHandler) HandleUserShow(c echo.Context) error {
	session, _ := session.Get("session", c)
	userSession := session.Values["user"].(goth.User)
	bodyBytes := clients.TopArtists(c)
	var topData model.SpotifyTopArtists
	err := json.Unmarshal(bodyBytes, &topData)
	if err != nil {
		log.Println("Unmarshall error")
	}
	var topArtists []model.SpotifyArtist
	for _, artist := range topData.Items {
		topArtists = append(topArtists, artist)

	}
	// spew.Dump(topArtists)
	userModel := model.User{
		Email:      userSession.Email,
		TopArtists: topArtists,
	}
	return render(c, user.Show(userModel))
}
