package handlers

import (
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
	topArtistsReponse := clients.TopArtists(userSession)
	var topArtists []model.SpotifyArtist
	for _, artist := range topArtistsReponse.Items {
		topArtists = append(topArtists, artist)

	}
	userModel := model.UserShowViewModel{
		Id:         userSession.UserID,
		Email:      userSession.Email,
		TopArtists: topArtists,
	}
	return render(c, user.Show(userModel))
}
