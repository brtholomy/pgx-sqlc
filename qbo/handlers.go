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

// FIXME: now I've got a struct wrapper passed to separate http.Handler implementations.
// Is this better?
type QboWrapper struct {
	client *qbohelp.Client
	// ErrorHandler func(r *http.Request, err error) http.Handler
}

func InitWrapper(c *qbohelp.Client) (QboWrapper, error) {
	if c == nil {
		return QboWrapper{}, errors.New("missing client!")
	}
	return QboWrapper{
		client: c,
	}, nil
}

type GetInvoiceHandler struct {
	wrapper QboWrapper
}

type PostInvoiceHandler struct {
	wrapper QboWrapper
}

func InitGetInvoiceHandler(w QboWrapper) GetInvoiceHandler {
	return GetInvoiceHandler{w}
}

func InitPostInvoiceHandler(w QboWrapper) PostInvoiceHandler {
	return PostInvoiceHandler{w}
}

func (h GetInvoiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleInvoice(w, r, nil)
}

func (h PostInvoiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: handle more than just amount:
	amount := ""
	if r.Form.Has("amount") {
		amount = r.Form.Get("amount")
	}

	invoice := fillInvoice(amount)
	resp, err := h.wrapper.client.CreateInvoice(&invoice)
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
