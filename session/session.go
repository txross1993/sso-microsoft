package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

func NewTokenFromOauth2Token(token *oauth2.Token) Token {
	return Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
}

func SetOauthState(w http.ResponseWriter) string {
	var expiration = time.Now().Add(24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	if encoded, err := cookieHandler.Encode("oauthstate", state); err == nil {
		cookie := &http.Cookie{
			Name:    "oauthstate",
			Value:   encoded,
			Expires: expiration,
		}
		http.SetCookie(w, cookie)
	}
	return state
}

func SetSession(token Token, response http.ResponseWriter) {
	value := map[string]string{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func GetOauthState(r *http.Request) string {
	var state string
	if cookie, err := r.Cookie("oauthstate"); err == nil {
		cookieHandler.Decode("oauthstate", cookie.Value, &state)
	}
	return state
}

func Clear(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func GetAccessToken(r *http.Request) *Token {
	var accessToken, refreshToken string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			accessToken = cookieValue["access_token"]
			refreshToken = cookieValue["refresh_token"]
		}
	}
	if accessToken == "" {
		return nil
	}
	return &Token{AccessToken: accessToken, RefreshToken: refreshToken}
}
