package handlers

import "github.com/labstack/echo/v4"

func SetupRoutes(e *echo.Echo, auth *AuthHandler, user *UserHandler, playlist *PlaylistHandler) {
	e.GET("/login", auth.HandleLoginShow)
	e.GET("/logout", auth.logoutHandler)
	e.GET("/auth/spotify", auth.spotifyLoginHandler)
	e.GET("/auth/spotify/callback", auth.spotifyCallbackHandler)
	e.GET("/auth/deezer", auth.deezerLoginHandler)
	e.GET("/auth/deezer/callback", auth.deezerCallbackHandler)

	protectedGroup := e.Group("/app", auth.authMiddleware)
	protectedGroup.GET("/user", user.HandleUserShow)
	protectedGroup.GET("/playlist", playlist.HandlePlaylistShow)
	protectedGroup.GET("/playlist/searchTracks", playlist.HandlePlaylistSearchTracks)
	protectedGroup.GET("/playlist/addTrackToPlaylist", playlist.HandleAddTrackToPlaylist)
}
