package db

import (
	"context"
	"log"
	"net/http"
	"pgx-sqlc/db/sqlc"
	"pgx-sqlc/ui/pages"

	"github.com/jackc/pgx/v5/pgtype"
)

const LOCALINV = "0197ada4-5f8b-77d7-b039-5651eabf1900"

// //////////////////////////////////////////////////////////
// DB

func fetchInvoiceItems(ctx context.Context, udb *UserDatabase, invid pgtype.UUID) ([]sqlc.InvoiceItem, error) {
	items, err := udb.DB.Sqlc.ListInvoiceItems(ctx, sqlc.ListInvoiceItemsParams{
		InvoiceID: invid,
		UserID:    udb.User.ID,
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

// //////////////////////////////////////////////////////////
// renderers

func renderInvoice(w http.ResponseWriter, r *http.Request, items []sqlc.InvoiceItem) {
	log.Println(items)
	component := pages.Debug(items)
	component.Render(r.Context(), w)
}

// //////////////////////////////////////////////////////////
// handlers

func GetInvoice(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {
	// TODO: take input
	invid, err := ReadUUID(LOCALINV)
	if err != nil {
		panic(err)
	}
	log.Printf("invid: %#v\n", invid)
	items, err := fetchInvoiceItems(ctx, dh.Udb, invid)
	if err != nil {
		panic(err)
	}
	renderInvoice(w, r, items)
}

// func PostInvoice(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	// TODO: handle more than just amount:
// 	name := ""
// 	if r.Form.Has("name") {
// 		name = r.Form.Get("name")
// 	}
// 	product, err := dh.Udb.NewInvoice(ctx, name)
// 	// TODO: handleError
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("new product: %#v\n", product)
// 	GetInvoice(ctx, dh, w, r)
// }
