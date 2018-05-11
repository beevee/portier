package yandex

import (
	"fmt"
	"net/url"
	"time"

	"github.com/levigross/grequests"
	"github.com/ssgreg/repeat"
)

const (
	apiURL        = "https://business.taxi.yandex.ru/api/1.0"
	passportURL   = "https://passport.yandex.ru/auth"
	loginAttempts = 10
)

type API struct {
	clientID string
	login    string
	password string
	session  *grequests.Session
}

func (a *API) Init(login, password, clientID string) error {
	if login == "" || password == "" || clientID == "" {
		return fmt.Errorf("login, password and clientID must all be specified")
	}

	a.login = login
	a.password = password
	a.clientID = clientID

	return repeat.Repeat(
		repeat.FnWithCounter(a.tryLogin),
		repeat.StopOnSuccess(),
		repeat.LimitMaxTries(loginAttempts),
		repeat.WithDelay(
			repeat.FullJitterBackoff(time.Minute).Set(),
		),
	)
}

func (a *API) tryLogin(attempt int) error {
	yandexURL, _ := url.Parse("https://yandex.ru")

	a.session = grequests.NewSession(&grequests.RequestOptions{
		UserAgent: "Mozilla/5.0 (Windows NT 6.1; rv:52.0) Gecko/20100101 Firefox/52.0",
		Headers: map[string]string{
			"Accept": "*/*",
		},
	})

	a.session.Post(passportURL, &grequests.RequestOptions{
		Data: map[string]string{
			"login":    a.login,
			"password": a.password,
		},
	})

	if len(a.session.HTTPClient.Jar.Cookies(yandexURL)) < 5 {
		return repeat.HintTemporary(fmt.Errorf("failed to log in for attempt %d", attempt))
	}

	return nil
}
