package qbo

import (
	"os"

	qbohelp "github.com/rwestlund/quickbooks-go"
)

// //////////////////////////////////////////////////////////
// QBO client

func loadClient(token *qbohelp.BearerToken) (*qbohelp.Client, error) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	realmId := os.Getenv("REALM_ID")
	// TODO: handle dev vs prod:
	return qbohelp.NewClient(clientId, clientSecret, realmId, false, "", token)
}

func SetupQboClient() (*qbohelp.Client, error) {
	// TODO: load from DB:
	bearer_token := &qbohelp.BearerToken{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
	}

	client, err := loadClient(bearer_token)
	if err != nil {
		return nil, err
	}

	// To do first when you receive the authorization code from quickbooks callback
	// authorizationCode := ""
	// redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"
	// bearerToken, err := client.RetrieveBearerToken(authorizationCode, redirectURI)
	// if err != nil {
	// 	panic(err)
	// }

	// TODO: figure out how often to refresh?
	bearer_token, err = client.RefreshToken(bearer_token.RefreshToken)
	if err != nil {
		return nil, err
	}
	return loadClient(bearer_token)
}
