package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"pgx-sqlc/ui/pages"

	qbo "github.com/rwestlund/quickbooks-go"
)

// //////////////////////////////////////////////////////////
// QBO client

func loadClient(token *qbo.BearerToken) (c *qbo.Client, err error) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	realmId := os.Getenv("REALM_ID")
	// TODO: handle dev vs prod:
	return qbo.NewClient(clientId, clientSecret, realmId, false, "", token)
}

func SetupQboClient() *qbo.Client {
	// FIXME: load from DB:
	token := qbo.BearerToken{
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

// //////////////////////////////////////////////////////////
// Invoice handling

func fillInvoice(amount string) qbo.Invoice {
	var invoice qbo.Invoice
	if err := json.Unmarshal([]byte(INVOICE), &invoice); err != nil {
		panic(err)
	}
	invoice.Line[0].Amount = json.Number(amount)
	return invoice
}

func handleInvoice(w http.ResponseWriter, r *http.Request, invoice *qbo.Invoice) {
	var resp pages.QboInvoiceResp
	if invoice != nil {
		jsonBytes, err := json.MarshalIndent(invoice, "", "  ")
		if err != nil {
			panic(err)
		}
		resp.InvoiceStr = string(jsonBytes)
		resp.Invoice = invoice
		resp.Amount = string(invoice.TotalAmt)
	}
	// spitting out component is not necessary, just for clarity:
	component := pages.Qbo(resp)
	component.Render(r.Context(), w)
}

// //////////////////////////////////////////////////////////
// Interface version
// http.Handler:
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

// FIXME: now I've got a struct wrapper passed to separate http.Handler implementations.
// Is this better?
type QboWrapper struct {
	Client *qbo.Client
	// ErrorHandler func(r *http.Request, err error) http.Handler
}

type GetQboHandler struct {
	Qbo *QboWrapper
}

type PostQboHandler struct {
	Qbo *QboWrapper
}

func (h GetQboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleInvoice(w, r, nil)
}

func (h PostQboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: handle more than just amount:
	amount := ""
	if r.Form.Has("amount") {
		amount = r.Form.Get("amount")
	}

	invoice := fillInvoice(amount)
	resp, err := h.Qbo.Client.CreateInvoice(&invoice)
	// TODO: what to do with errors? handleError()?
	if err != nil {
		panic(err)
	}
	handleInvoice(w, r, resp)
}

// //////////////////////////////////////////////////////////
// HandleFunc versions:

func QboGetHandler(w http.ResponseWriter, r *http.Request) {
	handleInvoice(w, r, nil)
}

func QboPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: handle more than just amount:
	amount := ""
	if r.Form.Has("amount") {
		amount = r.Form.Get("amount")
	}

	client := SetupQboClient()
	invoice := fillInvoice(amount)
	resp, err := client.CreateInvoice(&invoice)
	// TODO: what to do with errors? handleError()?
	if err != nil {
		panic(err)
	}
	handleInvoice(w, r, resp)
}
