package handlers

import (
	"strings"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/logger"
	"github.com/anthonymq/go-stack-demo/model"
	"github.com/anthonymq/go-stack-demo/view/playlist"
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
	var searchResults []model.SearchResult
	if userSession.Provider == "spotify" {
		searchResults = clients.SearchTrack(userSession, query).ToSearchResult()
	}
	if userSession.Provider == "deezer" {
		logger.Get().Info("deezer search")
		searchResults = clients.DeezerSearchTrack(userSession, query).ToSearchResults()
	}
	return (render(c, playlist.SearchResults(searchResults)))
}
