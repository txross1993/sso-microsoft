package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/txross1993/sso-microsoft/config"
	"github.com/txross1993/sso-microsoft/session"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

/*
	Resources
	https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-device-token-authentication
	  - Device Token for an interactive login session via a Microsoft login site

*/

var (
	router     = mux.NewRouter()
	AccessType = oauth2.AccessTypeOffline

	Endpoint = oauth2.Endpoint{
		AuthURL:   config.AuthURL,
		TokenURL:  config.TokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	}

	Config = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     Endpoint,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
	}
)

const (
	internalPage = `<h1>Internal</h1>
	<hr>
	<small>Access Token: %s</small>
	<form method="post" action="/logout">
		<button type="submit">Logout</button>
	</form>`

	index = `<h1>Login</h1>
	<form method="post" action="/auth/microsoft/login">
		<button type="submit">Login</button>
	</form>`
)

func main() {
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalHandler)
	router.HandleFunc("/auth/microsoft/login", loginHandler).Methods("POST")
	router.HandleFunc("/auth/microsoft/callback", oathRedirect)
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, index)
}

func internalHandler(w http.ResponseWriter, r *http.Request) {
	at := session.GetAccessToken(r)
	if at != nil {
		fmt.Fprintf(w, internalPage, at.AccessToken)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := session.SetOauthState(w)
	authURL := Config.AuthCodeURL(state,
		AccessType,
	)
	http.Redirect(w, r, authURL, 303)
}

func oathRedirect(w http.ResponseWriter, r *http.Request) {
	// read state from cookie
	oauthstate := session.GetOauthState(r)

	if r.FormValue("state") != oauthstate {
		log.Printf("invalid oauth state: GOT %s, WANT %s\n", r.FormValue("state"), oauthstate)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	authorizationCode := r.URL.Query().Get("code")
	token, err := Config.Exchange(
		context.Background(),
		authorizationCode,
		AccessType,
	)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`failed to aquire token: %s`, err.Error())))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	session.SetSession(session.NewTokenFromOauth2Token(token), w)
	http.Redirect(w, r, "/internal", 303)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session.Clear(w)
	http.Redirect(w, r, "/", 302)
}
