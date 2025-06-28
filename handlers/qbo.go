package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"pgx-sqlc/ui/pages"

	quickbooks "github.com/rwestlund/quickbooks-go"
)

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

// TODO: is a http.Handler type preferable because it could store the qbo client and other reusables?

// // http.Handler:
// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }

// func NewQboHandler(process func(*http.Request) QboInvoiceResp) QboHandler {
// 	return QboHandler{Process: process}
// }

// type QboHandler struct {
// 	// gets called within the ServeHTTP() method.
// 	Process func(r *http.Request) QboInvoiceResp
// }

// // implements the HTTP handler interface on the QboHandler type.
// func (qh QboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// NOTE: can pass in http.Request fields and methods.
// 	// TODO: change this signature from returning a string.
// 	resp := qh.Process(r)
// 	pages.Qbo(resp.Amount, resp.Resp).Render(r.Context(), w)
// }

// func GetInvoice(r *http.Request, amount string, invoice *quickbooks.Invoice) QboInvoiceResp {
// 	var resp QboInvoiceResp
// 	if invoice != nil {
// 		jsonBytes, err := json.MarshalIndent(invoice, "", "  ")
// 		if err != nil {
// 			panic(err)
// 		}
// 		invstr = string(jsonBytes)
// 	}
// 	return resp
// }

// //////////////////////////////////////////////////////////
// HandleFunc versions:
func handleInvoice(w http.ResponseWriter, r *http.Request, invoice *quickbooks.Invoice) {
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

	client := setupQboClient()
	invoice := fillInvoice(amount)
	resp, err := client.CreateInvoice(&invoice)
	// TODO: what to do with errors? handleError()?
	if err != nil {
		panic(err)
	}
	handleInvoice(w, r, resp)
}
