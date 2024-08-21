package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/anthonymq/go-stack-demo/logger"
	"github.com/anthonymq/go-stack-demo/model"
	"github.com/davecgh/go-spew/spew"
	"github.com/markbates/goth"
	"go.uber.org/zap"
)

const (
	authURL         string = "https://connect.deezer.com/oauth/auth.php"
	tokenURL        string = "https://connect.deezer.com/oauth/access_token.php"
	endpointProfile string = "https://api.deezer.com/user/me"
	apiBasePath     string = "https://api.deezer.com"
)

// DeezerSession stores data during the auth process with Deezer.
type DeezerSession struct {
	AuthURL     string
	AccessToken string
	ExpiresAt   time.Time
}

var DeezerProvider goth.Provider

func FetchUser(session DeezerSession) (goth.User, error) {
	sess := session
	user := goth.User{
		AccessToken: session.AccessToken,
		Provider:    "deezer",
		ExpiresAt:   session.ExpiresAt,
	}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", "deezer")
	}
	response, err := http.Get(endpointProfile + "?access_token=" + url.QueryEscape(sess.AccessToken))
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("%s responded with a %d trying to fetch user information", "deezer", response.StatusCode)
	}

	bits, err := io.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(bits)).Decode(&user.RawData)
	if err != nil {
		return user, err
	}

	err = userFromReader(bytes.NewReader(bits), &user)
	return user, err

}

func userFromReader(reader io.Reader, user *goth.User) error {
	u := struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		NickName  string `json:"name"`
		AvatarURL string `json:"picture"`
		Location  string `json:"city"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		spew.Dump(err)
		return err
	}

	user.UserID = strconv.Itoa(u.ID)
	user.Email = u.Email
	user.FirstName = u.FirstName
	user.LastName = u.LastName
	user.NickName = u.NickName
	user.AvatarURL = u.AvatarURL
	user.Location = u.Location

	return nil
}

func callDeezerApi(path string, u goth.User, urlParams url.Values) (*http.Response, error) {
	client := &http.Client{}
	base, err := url.Parse(apiBasePath + path)
	if err != nil {
		return nil, err
	}
	base.RawQuery = urlParams.Encode()
	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		logger.Get().Error(err.Error())
		return nil, err
	}
	return client.Do(req)
}
func DeezerSearchTrack(u goth.User, query string) model.DeezerSearchTrackResults {
	params := url.Values{}
	params.Add("q", query)
	params.Add("limit", "10")
	resp, err := callDeezerApi(
		"/search",
		u,
		params,
	)
	defer resp.Body.Close()
	if err != nil {
		logger.Get().Error(err.Error())
	}
	logger.Get().Info(resp.Status)
	results, err := unmarshal[model.DeezerSearchTrackResults](resp)
	if err != nil {
		logger.Get().Error(err.Error())
	}
	return results
}

func DeezerGetUserPlaylists(u goth.User) model.DeezerGetPlaylists {
	resp, err := callDeezerApi(
		fmt.Sprintf("/user/%s/playlists", u.UserID),
		u,
		nil,
	)

	defer resp.Body.Close()
	if err != nil {
		logger.Get().Error(err.Error())
	}
	logger.Get().Info(resp.Status)
	results, err := unmarshal[model.DeezerGetPlaylists](resp)
	if err != nil {
		logger.Get().Error(err.Error())
	}
	return results

}

func DeezerListeningHistory(u goth.User) {
	resp, err := callDeezerApi("/user/me/history", u, nil)
	defer resp.Body.Close()

	if err != nil {
		logger.Get().Error(err.Error())
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	logger.Get().Debug("Response body", zap.String("body", string(bodyBytes)))
	logger.Get().Info(resp.Status)

}

func DeezerAddTrackToPlaylist(u goth.User, plId int, tId string) bool {
	path := fmt.Sprintf("/playlist/%d/tracks", plId)
	params := url.Values{}
	params.Add("request_method", "POST")
	params.Add("songs", tId)
	// params.Add("order", tId)
	params.Add("access_token", u.AccessToken)
	resp, err := callDeezerApi(
		path,
		u,
		params,
	)

	defer resp.Body.Close()
	if err != nil {
		logger.Get().Error(err.Error())
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	logger.Get().Debug("Response body", zap.String("body", string(bodyBytes)))
	logger.Get().Info(resp.Status)
	return resp.StatusCode == http.StatusOK
}
