package db

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pgx-sqlc/db/sqlc"
	"pgx-sqlc/ui/pages"

	"github.com/jackc/pgx/v5/pgtype"
)

// //////////////////////////////////////////////////////////
// DB

func newProduct(ctx context.Context, udb *UserDatabase, name, price string) (*sqlc.Product, error) {
	pid, err := MakeUUIDv7()
	if err != nil {
		return nil, err
	}
	var num pgtype.Numeric
	err = num.Scan(price)
	if err != nil {
		return nil, err
	}
	newprod, err := udb.DB.Sqlc.CreateProduct(ctx, sqlc.CreateProductParams{
		ID:     pid,
		UserID: udb.User.ID,
		Name:   name,
		Price:  num,
	})
	if err != nil {
		return nil, err
	}
	return &newprod, nil
}

func listProducts(ctx context.Context, udb *UserDatabase) ([]sqlc.Product, error) {
	products, err := udb.DB.Sqlc.ListProducts(ctx, udb.User.ID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func convertToPageProduct(p sqlc.Product) (pages.Product, error) {
	// TODO: is float64 what we want? Seems to be the only option from pgtypes.Numeric.
	pgtype_int, err := p.Price.Float64Value()
	if err != nil {
		log.Printf("Failed to parse price for: %#v. err: %v\n", p.Name, err)
		return pages.Product{}, err
	}
	// https://github.com/a-h/templ/issues/307#issuecomment-1828720574
	price := fmt.Sprintf("%.2f", pgtype_int.Float64)
	pidstr := p.ID.String()
	return pages.Product{p.Name, price, pidstr}, nil
}

func convertToPageProducts(sps []sqlc.Product) []pages.Product {
	var products []pages.Product
	for _, sp := range sps {
		p, err := convertToPageProduct(sp)
		if err != nil {
			log.Printf("failed to convert product: %#v\n", err)
			continue
		}
		products = append(products, p)
	}
	return products
}

// //////////////////////////////////////////////////////////
// page renderers

func renderProducts(w http.ResponseWriter, r *http.Request, products []pages.Product) {
	// TODO: do something when products is empty.
	component := pages.ListProducts(products)
	component.Render(r.Context(), w)
}

// //////////////////////////////////////////////////////////
// handlers

func GetProducts(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {
	sps, err := listProducts(ctx, dh.Udb)
	if err != nil {
		panic(err)
	}
	products := convertToPageProducts(sps)
	renderProducts(w, r, products)
}

func PostProducts(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: handle more than just amount:
	name := ""
	price := ""
	if r.Form.Has("name") {
		name = r.Form.Get("name")
	}
	if r.Form.Has("price") {
		price = r.Form.Get("price")
	}
	product, err := newProduct(ctx, dh.Udb, name, price)
	// TODO: handleError
	if err != nil {
		panic(err)
	}
	log.Printf("new product: %#v\n", product)
	GetProducts(ctx, dh, w, r)
}
