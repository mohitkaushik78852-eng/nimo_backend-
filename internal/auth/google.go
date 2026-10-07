package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GoogleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Aud           string `json:"aud"`
	ExpiresIn     string `json:"expires_in"`
}

func VerifyGoogleIDToken(idToken string, expectedClientID string) (*GoogleTokenInfo, error) {
	if expectedClientID == "" {
		return nil, errors.New("google sign-in is not configured")
	}

	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(endpoint)
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

	if info.Aud != expectedClientID {
		return nil, errors.New("google id token was not issued for this app")
	}
	if info.Sub == "" || info.Email == "" || info.EmailVerified != "true" {
		return nil, errors.New("google email is not verified")
	}
	if expiresIn, err := strconv.Atoi(info.ExpiresIn); err != nil || expiresIn <= 0 {
		return nil, errors.New("google id token has expired")
	}

	return &info, nil
}
