package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

/*
Basic Authorization URLS

https://login.microsoftonline.com/common/oauth2/v2.0/authorize
https://login.microsoftonline.com/common/oauth2/v2.0/token

To ensure compatibility of your app in Safari and other privacy-conscious browsers, we no longer recommend use of the implicit flow and instead recommend the authorization code flow.
*/

var (
	TenantID, ClientID, ClientSecret, Host, Port, RedirectURL, msAuthBase, msTokenBase, BaseURL, AuthURL, TokenURL string
	Scopes                                                                                                         []string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	TenantID = os.Getenv("TENANT_ID")
	ClientID = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	Host = os.Getenv("HOST")
	Port = os.Getenv("PORT")
	RedirectURL = os.Getenv("REDIRECT_URL")

	msAuthBase = "https://login.microsoftonline.com/%s/oauth2/v2.0/authorize"
	msTokenBase = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	BaseURL = fmt.Sprintf("http://%s:%s", Host, Port)
	Scopes = []string{"openid", "https://graph.microsoft.com/user.read"}
	AuthURL = fmt.Sprintf(msAuthBase, TenantID)
	TokenURL = fmt.Sprintf(msTokenBase, TenantID)
}

func ScopesAsString() string {
	var scope string
	for _, s := range Scopes {
		scope += fmt.Sprintf("%s ", s)
	}
	return strings.TrimRight(scope, " ")
}
