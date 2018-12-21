package yandex

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/levigross/grequests"
)

const (
	apiURL = "https://business.taxi.yandex.ru/api/1.0"
)

// API is a singleton client for Yandex.Taxi API
type API struct {
	clientID string
	session  *grequests.Session
}

// Init establishes authenticated session with Yandex.Passport
func (a *API) Init(sessionID, clientID string) error {
	if sessionID == "" || clientID == "" {
		return fmt.Errorf("sessionID and clientID must be specified")
	}

	a.clientID = clientID

	cookieJar, _ := cookiejar.New(nil)
	a.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: "Mozilla/5.0 (Windows NT 6.1; rv:52.0) Gecko/20100101 Firefox/52.0",
		Headers: map[string]string{
			"Accept": "*/*",
		},
		CookieJar: cookieJar,
	})

	cookieUrl, _ := url.Parse(apiURL)
	a.session.RequestOptions.CookieJar.SetCookies(cookieUrl, []*http.Cookie{
		{
			Name:     "Session_id",
			Value:    sessionID,
			HttpOnly: true,
			Secure:   true,
		},
	})

	return nil
}
