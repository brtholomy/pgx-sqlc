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

// Fetches all InvoiceItems and their Products matching user and invoice id.
//
// First get matching InvoiceItems, then get matching Product for each.
func fetchInvoiceItems(
	ctx context.Context,
	udb *UserDatabase,
	invid pgtype.UUID,
) ([]pages.InvoiceItem, error) {
	sqlc_items, err := udb.DB.Sqlc.ListInvoiceItems(ctx, sqlc.ListInvoiceItemsParams{
		InvoiceID: invid,
		UserID:    udb.User.ID,
	})
	if err != nil {
		return nil, err
	}

	// NOTE: would rather do the conversion at the render step, but I need a wrapper for each item
	// anyway, since the sqlc.InvoiceItem does not contain the Product, just the ID.
	items := []pages.InvoiceItem{}
	for _, si := range sqlc_items {
		sp, err := udb.DB.Sqlc.GetProduct(ctx, si.ProductID)
		if err != nil {
			log.Printf("failed to get product: %#v\n", err)
			continue
		}
		p, err := convertToPageProduct(sp)
		if err != nil {
			log.Printf("failed to convert product: %#v\n", err)
			continue
		}
		i := pages.InvoiceItem{
			Product: p,
		}
		items = append(items, i)
	}

	return items, nil
}

// //////////////////////////////////////////////////////////
// renderers

func renderInvoice(w http.ResponseWriter, r *http.Request, items []pages.InvoiceItem) {
	log.Println(items)
	component := pages.Invoice(items)
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
