package qbo

import (
	"encoding/json"
	"errors"
	"net/http"

	"pgx-sqlc/ui/pages"

	qbohelp "github.com/rwestlund/quickbooks-go"
)

// //////////////////////////////////////////////////////////
// Invoice handling

func fillInvoice(amount string) qbohelp.Invoice {
	var invoice qbohelp.Invoice
	if err := json.Unmarshal([]byte(INVOICE), &invoice); err != nil {
		panic(err)
	}
	invoice.Line[0].Amount = json.Number(amount)
	return invoice
}

func handleInvoice(w http.ResponseWriter, r *http.Request, invoice *qbohelp.Invoice) {
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

// http.Handler that takes a function as dependency. See InitHandler()
type QboHandler struct {
	process func(qh *QboHandler, w http.ResponseWriter, r *http.Request)
	// Public because the process function needs it.
	Client *qbohelp.Client
}

// implements the HTTP handler interface on the QboHandler type.
func (qh QboHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: passing itself as parameter is essentially what a method receiver does:
	qh.process(&qh, w, r)
}

// Sets up a QboHandler while checking dependencies.
func InitHandler(
	c *qbohelp.Client,
	p func(qh *QboHandler, w http.ResponseWriter, r *http.Request)) (QboHandler, error) {
	if c == nil {
		return QboHandler{}, errors.New("missing client!")
	}
	if p == nil {
		return QboHandler{}, errors.New("missing function!")
	}
	return QboHandler{
		Client:  c,
		process: p,
	}, nil
}

func GetInvoice(qh *QboHandler, w http.ResponseWriter, r *http.Request) {
	handleInvoice(w, r, nil)
}

func PostInvoice(qh *QboHandler, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: handle more than just amount:
	amount := ""
	if r.Form.Has("amount") {
		amount = r.Form.Get("amount")
	}

	invoice := fillInvoice(amount)
	resp, err := qh.Client.CreateInvoice(&invoice)
	// TODO: what to do with errors? handleError()?
	if err != nil {
		panic(err)
	}
	handleInvoice(w, r, resp)
}

// //////////////////////////////////////////////////////////
// HandleFunc versions:

func GetInvoiceFunc(w http.ResponseWriter, r *http.Request) {
	handleInvoice(w, r, nil)
}

func PostInvoiceFunc(w http.ResponseWriter, r *http.Request) {
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
