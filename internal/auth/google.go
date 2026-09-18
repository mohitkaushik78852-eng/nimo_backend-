package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type GoogleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Aud           string `json:"aud"`
}

func VerifyGoogleIDToken(idToken string, expectedClientID string) (*GoogleTokenInfo, error) {
	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid google id token")
	}

	var info GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	if expectedClientID != "" && info.Aud != expectedClientID {
		return nil, errors.New("google id token was not issued for this app")
	}
	if info.EmailVerified != "true" {
		return nil, errors.New("google email is not verified")
	}

	return &info, nil
}
