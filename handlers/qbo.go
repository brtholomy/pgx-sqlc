package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"pgx-sqlc/ui/pages"

	quickbooks "github.com/rwestlund/quickbooks-go"
)

// TODO: figure out if this factory method thing is better.
// func NewQboHandler(request func(*http.Request) string) QboHandler {
// 	return QboHandler{Request: request}
// }

// type QboHandler struct {
// 	// must match the signature of the QboRequest func below.
// 	Request func(r *http.Request) string
// }

// // implements the HTTP handler interface on the QboHandler type.
// func (qh QboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// NOTE: can pass in http.Request fields and methods.
// 	// TODO: change this signature from returning a string.
// 	resp := qh.Request(r)
// 	pages.Qbo(resp).Render(r.Context(), w)
// }

func loadClient(token *quickbooks.BearerToken) (c *quickbooks.Client, err error) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	realmId := os.Getenv("REALM_ID")
	return quickbooks.NewClient(clientId, clientSecret, realmId, false, "", token)
}

func setupQboClient() *quickbooks.Client {
	// FIXME: load from DB:
	token := quickbooks.BearerToken{
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

func fillInvoice(amount string) quickbooks.Invoice {
	var invoice quickbooks.Invoice
	if err := json.Unmarshal([]byte(INVOICE), &invoice); err != nil {
		panic(err)
	}
	invoice.Line[0].Amount = json.Number(amount)
	return invoice
}

func QboPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	amount := ""
	if r.Form.Has("amount") {
		amount = r.Form.Get("amount")
	}

	client := setupQboClient()
	invoice := fillInvoice(amount)
	resp, err := client.CreateInvoice(&invoice)
	if err != nil {
		panic(err)
	}

	QboGetHandler(w, r, amount, resp)
}

func QboGetHandler(w http.ResponseWriter, r *http.Request, amount string, invoice *quickbooks.Invoice) {
	invstr := ""
	if invoice != nil {
		jsonBytes, err := json.MarshalIndent(invoice, "", "  ")
		if err != nil {
			panic(err)
		}
		invstr = string(jsonBytes)
	}
	component := pages.Qbo(amount, invstr)
	component.Render(r.Context(), w)
}
