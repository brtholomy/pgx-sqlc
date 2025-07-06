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

func newInvoiceItem(ctx context.Context, udb *UserDatabase, invid pgtype.UUID, pid pgtype.UUID) (*sqlc.InvoiceItem, error) {
	itemid, err := MakeUUIDv7()
	if err != nil {
		return nil, err
	}
	var num pgtype.Numeric
	// TODO: pass in product.price ? do a lookup?
	err = num.Scan("0.00")
	if err != nil {
		return nil, err
	}
	newprod, err := udb.DB.Sqlc.CreateInvoiceItem(ctx, sqlc.CreateInvoiceItemParams{
		ID:        itemid,
		UserID:    udb.User.ID,
		ProductID: pid,
		InvoiceID: invid,
		Amount:    num,
	})
	if err != nil {
		return nil, err
	}
	return &newprod, nil
}

// //////////////////////////////////////////////////////////
// renderers

func renderInvoice(w http.ResponseWriter, r *http.Request, i []pages.InvoiceItem, p []pages.Product) {
	log.Println(i)
	component := pages.Invoice(i, p)
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
	sps, err := listProducts(ctx, dh.Udb)
	if err != nil {
		panic(err)
	}
	products := convertToPageProducts(sps)
	renderInvoice(w, r, items, products)
}

func PostInvoice(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	pidstr := ""
	if r.Form.Has("product") {
		pidstr = r.Form.Get("product")
	}
	pid, err := ReadUUID(pidstr)
	if err != nil {
		panic(err)
	}
	// TODO: get back from req?
	invid, err := ReadUUID(LOCALINV)
	if err != nil {
		panic(err)
	}
	product, err := newInvoiceItem(ctx, dh.Udb, invid, pid)
	if err != nil {
		panic(err)
	}
	log.Printf("new product: %#v\n", product)
	GetInvoice(ctx, dh, w, r)
}
