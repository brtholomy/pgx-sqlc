package db

import (
	"context"
	"errors"
	"net/http"
)

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
