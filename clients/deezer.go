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
)

const (
	authURL         string = "https://connect.deezer.com/oauth/auth.php"
	tokenURL        string = "https://connect.deezer.com/oauth/access_token.php"
	endpointProfile string = "https://api.deezer.com/user/me"
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

func callDeezerApi(method string, url string, u goth.User) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		logger.Get().Error(err.Error())
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+u.AccessToken)
	return client.Do(req)
}
func DeezerSearchTrack(u goth.User, query string) model.DeezerSearchTrackResults {
	resp, err := callDeezerApi("GET",
		fmt.Sprintf("https://api.deezer.com/search?q=%s&limit=10", query),
		u,
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
	spew.Dump(results)
	return results
}
