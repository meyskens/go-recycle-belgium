package recyclebelgium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type API struct {
	authTokenMutex sync.Mutex
	authToken      string
	xsecret        string
}

func NewAPI(xsecert string) *API {
	return &API{
		xsecret: xsecert,
	}
}

type authToken struct {
	ExpiresAt   time.Time `json:"expiresAt"`
	AccessToken string    `json:"accessToken"`
}

func (a *API) getToken() (string, error) {
	a.authTokenMutex.Lock()
	defer a.authTokenMutex.Unlock()

	if a.authToken == "" {
		err := a.fetchAuthToken()
		if err != nil {
			return "", err
		}
	}

	return a.authToken, nil
}

func (a *API) fetchAuthToken() error {
	req, err := http.NewRequest("GET", "https://api.fostplus.be/recycle-public/app/v1/access-token", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Consumer", "recycleapp.be") // shoud be recycleapp.be or mobile-app, for some reason it wants only those two
	req.Header.Set("X-Secret", a.xsecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// fetch body
	var authToken authToken

	err = json.NewDecoder(resp.Body).Decode(&authToken)
	if err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	a.authToken = authToken.AccessToken

	// auto expire the token
	go func() {
		time.Sleep(time.Until(authToken.ExpiresAt))
		a.authTokenMutex.Lock()
		a.authToken = ""
		a.authTokenMutex.Unlock()
	}()

	return nil
}
