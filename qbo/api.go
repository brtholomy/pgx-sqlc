package qbo

import (
	"os"

	qbohelp "github.com/rwestlund/quickbooks-go"
)

// //////////////////////////////////////////////////////////
// QBO client

func loadClient(token *qbohelp.BearerToken) (c *qbohelp.Client, err error) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	realmId := os.Getenv("REALM_ID")
	// TODO: handle dev vs prod:
	return qbohelp.NewClient(clientId, clientSecret, realmId, false, "", token)
}

func SetupQboClient() *qbohelp.Client {
	// FIXME: load from DB:
	token := qbohelp.BearerToken{
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
	}

	client, err := loadClient(&token)
	if err != nil {
		panic(err)
	}

	// To do first when you receive the authorization code from quickbooks callback
	// authorizationCode := "XAB11746551225hXNdSW2iGUcTdTLImx5gzNIF59QnhMmM40tX"
	// redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"
	// bearerToken, err := client.RetrieveBearerToken(authorizationCode, redirectURI)
	// if err != nil {
	// 	panic(err)
	// }

	// TODO: figure out how often to refresh?
	_, err = client.RefreshToken(token.RefreshToken)
	if err != nil {
		panic(err)
	}
	return client
}
