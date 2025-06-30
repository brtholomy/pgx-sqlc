package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"pgx-sqlc/ui/pages"

	qbo "github.com/rwestlund/quickbooks-go"
)

func loadClient(token *qbo.BearerToken) (c *qbo.Client, err error) {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("SECRET")
	realmId := os.Getenv("REALM_ID")
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

func fillInvoice(amount string) qbo.Invoice {
	var invoice qbo.Invoice
	if err := json.Unmarshal([]byte(INVOICE), &invoice); err != nil {
		panic(err)
	}
	invoice.Line[0].Amount = json.Number(amount)
	return invoice
}

// //////////////////////////////////////////////////////////
// Interface version
// http.Handler:
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

// FIXME: this pattern seems wrong.
// Passing the client to the constructor is correct.
// Passing the function to the constructor is correct.
// But passing the client through to the function seems wrong? I wanted to use the struct client
// field from within the qh.QboPostHandler method. But how would I pass a method of an object to its
// own constructor?
//
// Or do I need a type for each endpoint? That's not very DRY.
// QboGetHandler
// NewQboGetHandler(client)
type QboHandler struct {
	Process func(w http.ResponseWriter, r *http.Request, c *qbo.Client)
	Client  *qbo.Client
}

// implements the HTTP handler interface on the QboHandler type.
func (qh QboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	qh.Process(w, r, qh.Client)
}

// NOTE: this isn't a method so it can be passed to the constructor:
func GetHandleInvoice(w http.ResponseWriter, r *http.Request, c *qbo.Client) {
	handleInvoice(w, r, nil)
}

func PostHandleInvoice(w http.ResponseWriter, r *http.Request, c *qbo.Client) {
	r.ParseForm()

	// TODO: handle more than just amount:
	amount := ""
	if r.Form.Has("amount") {
		amount = r.Form.Get("amount")
	}

	invoice := fillInvoice(amount)
	resp, err := c.CreateInvoice(&invoice)
	// TODO: what to do with errors? handleError()?
	if err != nil {
		panic(err)
	}
	handleInvoice(w, r, resp)
}

// //////////////////////////////////////////////////////////
// HandleFunc versions:
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
