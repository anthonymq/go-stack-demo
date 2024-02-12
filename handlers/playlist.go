package handlers

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/model"
	"github.com/anthonymq/go-stack-demo/view/playlist"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
)

type PlaylistHandler struct {
}

func (h PlaylistHandler) HandlePlaylistShow(c echo.Context) error {
	return (render(c, playlist.Show()))
}

func (h PlaylistHandler) HandlePlaylistSearchTracks(c echo.Context) error {
	session, _ := session.Get("session", c)
	userSession := session.Values["user"].(goth.User)
	query := strings.Replace(c.QueryParam("search"), " ", "+", -1)
	bodyBytes := clients.SearchTrack(userSession, query)
	var searchResults model.SearchResult
	err := json.Unmarshal(bodyBytes, &searchResults)
	if err != nil {
		spew.Dump(err)
		log.Println("Unmarshall error")
	}

	return (render(c, playlist.SearchResults(searchResults)))
}
