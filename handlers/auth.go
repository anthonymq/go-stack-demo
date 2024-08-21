package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/anthonymq/go-stack-demo/clients"
	"github.com/anthonymq/go-stack-demo/logger"
	"github.com/anthonymq/go-stack-demo/view"
	"github.com/anthonymq/go-stack-demo/view/login"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct{}

func (h *AuthHandler) logoutHandler(c echo.Context) error {
	// Get the session from the request
	sess, err := session.Get("session", c)
	if err != nil {
		logger.Get().Error(err.Error())
	}
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (h AuthHandler) HandleLoginShow(c echo.Context) error {
	return render(c, login.Show())
}

func (h *AuthHandler) spotifyLoginHandler(c echo.Context) error {
	gothic.GetProviderName = func(req *http.Request) (string, error) { return "spotify", nil }
	if _, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
		logger.Get().Info("Already loggedIn")
		h.HandleLoginShow(c)
	} else {
		logger.Get().Info("Not already loggedIn")
		gothic.BeginAuthHandler(c.Response(), c.Request())
	}
	return nil
}
func (h *AuthHandler) deezerLoginHandler(c echo.Context) error {
	gothic.GetProviderName = func(req *http.Request) (string, error) { return "deezer", nil }
	if _, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
		logger.Get().Info("Already loggedIn")
		h.HandleLoginShow(c)
	} else {
		logger.Get().Info("Not already loggedIn")
		gothic.BeginAuthHandler(c.Response(), c.Request())
	}
	return nil
}

func (h *AuthHandler) spotifyCallbackHandler(c echo.Context) error {
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
}

func (h *AuthHandler) deezerCallbackHandler(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	session, err := session.Get("session", c)
	if err != nil {
		log.Println("Error retrieving session", err)
	}
	user.RawData = map[string]interface{}{}
	session.Values["user"] = user
	spew.Dump(user)
	err = session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		log.Println("Error saving session", err)
	}
	return c.Redirect(http.StatusSeeOther, "/app/user")
}

func (h *AuthHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		userConnected, ok := GetUserInContext(c)
		if !ok {
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "please provide valid credentials")
		}

		if userConnected.Provider == "spotify" {
			clients.RenewAccessTokenIfExpired(c, *sess, userConnected)
		}

		var avatarUrl string
		avatarUrl = userConnected.AvatarURL

		ctx := context.WithValue(c.Request().Context(), view.ContextAvatarUrlKey, avatarUrl)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
func PopulateTemplContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var avatarUrl string
		userConnected, ok := GetUserInContext(c)
		if ok {
			avatarUrl = userConnected.AvatarURL
		}
		ctx := context.WithValue(c.Request().Context(), view.ContextAvatarUrlKey, avatarUrl)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func GetUserInContext(c echo.Context) (goth.User, bool) {
	sess, _ := session.Get("session", c)
	userConnected, ok := sess.Values["user"].(goth.User)
	return userConnected, ok
}
