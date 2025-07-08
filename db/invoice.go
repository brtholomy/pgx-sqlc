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
			ID:      si.ID.String(),
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

func deleteInvoiceItem(ctx context.Context, udb *UserDatabase, item_id pgtype.UUID) error {
	err := udb.DB.Sqlc.DeleteInvoiceItem(ctx, item_id)
	if err != nil {
		return err
	}
	return nil
}

// //////////////////////////////////////////////////////////
// renderers

func renderInvoice(w http.ResponseWriter, r *http.Request, i []pages.InvoiceItem, p []pages.Product) {
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

	// TODO: get back from req?
	invid, err := ReadUUID(LOCALINV)
	if err != nil {
		panic(err)
	}

	pid, err := GetUUIDFromUrlValues(r.Form, "add-product")
	if err != nil {
		log.Printf("err: %#v\n", err)
		log.Printf("req: %#v\n", r)
		pages.Debug(r).Render(r.Context(), w)
		return
	}
	invitem, err := newInvoiceItem(ctx, dh.Udb, invid, pid)
	if err != nil {
		panic(err)
	}
	log.Printf("new product: %#v\n", invitem.ProductID.String())

	items, err := fetchInvoiceItems(ctx, dh.Udb, invid)
	if err != nil {
		panic(err)
	}
	pages.DisplayInvoice(items).Render(r.Context(), w)
}

func DeleteInvoiceItem(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {

	// TODO: get back from req?
	invid, err := ReadUUID(LOCALINV)
	if err != nil {
		panic(err)
	}

	item_id, err := GetUUIDFromUrlValues(r.URL.Query(), "delete-invoice-item")
	if err != nil {
		log.Printf("err: %#v\n", err)
		log.Printf("req: %#v\n", r)
		pages.Debug(r).Render(r.Context(), w)
		return
	}
	err = deleteInvoiceItem(ctx, dh.Udb, item_id)
	if err != nil {
		log.Fatalf("delete failed. item_id: %v", item_id, err)
		pages.Debug(r).Render(r.Context(), w)
		return
	}
	log.Printf("deleted invoice_item: %#v\n", item_id.String())

	items, err := fetchInvoiceItems(ctx, dh.Udb, invid)
	if err != nil {
		panic(err)
	}
	pages.DisplayInvoice(items).Render(r.Context(), w)
}
