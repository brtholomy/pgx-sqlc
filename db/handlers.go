package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"pgx-sqlc/db/sqlc"
	"pgx-sqlc/ui/pages"
)

// //////////////////////////////////////////////////////////
// helpers

func handleProducts(w http.ResponseWriter, r *http.Request, in []sqlc.Product) {
	// func handleProducts(w http.ResponseWriter, r *http.Request, products []sqlc.Product) {

	var products []pages.Product
	// FIXME: do something when products is empty.
	for _, p := range in {
		// TODO: is float64 what we want? Seems to be the only option from pgtypes.Numeric that
		// works.
		pgtype_int, err := p.Price.Float64Value()
		if err != nil {
			log.Printf("Failed to parse price for: %#v. err: %v\n", p.Name, err)
		} else {
			// https://github.com/a-h/templ/issues/307#issuecomment-1828720574
			price := fmt.Sprintf("%f", pgtype_int.Float64)
			products = append(products, pages.Product{p.Name, price})
		}
	}

	component := pages.ListProducts(products)
	component.Render(r.Context(), w)
}

// //////////////////////////////////////////////////////////
// http.Handler

// http.Handler that takes a function as dependency. See InitHandler()
type DbHandler struct {
	process func(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request)
	// Public because the process function needs it.
	Udb *UserDatabase
}

// implements the HTTP handler interface on the DbHandler type.
func (dh DbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: passing itself as parameter is essentially what a method receiver does:
	ctx := context.Background()
	dh.process(ctx, &dh, w, r)
}

// Sets up a DbHandler while checking dependencies.
func InitHandler(
	udb *UserDatabase,
	p func(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request)) (DbHandler, error) {
	if udb == nil {
		return DbHandler{}, errors.New("missing client!")
	}
	if p == nil {
		return DbHandler{}, errors.New("missing function!")
	}
	return DbHandler{
		Udb:     udb,
		process: p,
	}, nil
}

func GetProducts(ctx context.Context, dh *DbHandler, w http.ResponseWriter, r *http.Request) {
	products, err := dh.Udb.ListProducts(ctx)
	if err != nil {
		panic(err)
	}
	// log.Println(products)
	handleProducts(w, r, products)
}

// func PostProducts(dh *DbHandler, w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	// TODO: handle more than just amount:
// 	amount := ""
// 	if r.Form.Has("amount") {
// 		amount = r.Form.Get("amount")
// 	}

// 	invoice := fillProducts(amount)
// 	resp, err := dh.udb.DB.Query.CreateProduct()
// 	// TODO: what to do with errors? handleError()?
// 	if err != nil {
// 		panic(err)
// 	}
// 	handleProducts(w, r, resp)
// }
