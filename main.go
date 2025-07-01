package main

import (
	"fmt"
	"net/http"
	"os"

	"pgx-sqlc/qbo"
	"pgx-sqlc/ui/assets"
	"pgx-sqlc/ui/pages"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
)

func initDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func setupAssetsRoutes(mux *http.ServeMux) {
	var is_development = os.Getenv("GO_ENV") != "production"

	asset_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if is_development {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if is_development {
			fs = http.FileServer(http.Dir("./ui/assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /ui/assets/", http.StripPrefix("/ui/assets/", asset_handler))
}

func main() {
	initDotEnv()
	mux := http.NewServeMux()
	setupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))

	// http.Handler implementations:
	wrapper := qbo.QboWrapper{qbo.SetupQboClient()}
	mux.Handle("GET /qbo", qbo.GetInvoiceHandler{wrapper})
	mux.Handle("POST /qbo", qbo.PostInvoiceHandler{wrapper})

	// HandleFunc versions:
	// mux.HandleFunc("GET /qbo", qbo.GetInvoiceFunc)
	// mux.HandleFunc("POST /qbo", qbo.PostInvoiceFunc)

	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)

	// DB code:
	// ctx := context.Background()
	// db_url := "user=bth database=testdb"
	// pgdb := db.NewDatabase(ctx, db_url)

	// // joe, err := db.NewUser(ctx, pgdb, "joe", "j@blow.com")
	// joe, err := db.GetUser(ctx, pgdb, "0197ada4-5f8b-77d7-b039-5651eabf19e1")
	// if err != nil {
	// 	panic(err)
	// }
	// udb := db.UserDatabase{&joe, pgdb}
	// newprod, err := udb.NewProduct(ctx, "stuff", "3.15")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(newprod.Name)

	// products, err := pgdb.Query.ListProducts(ctx, joe.ID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(products)
}
