package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"pgx-sqlc/db"
	"pgx-sqlc/ui/assets"
	"pgx-sqlc/ui/pages"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
)

func InitDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func SetupAssetsRoutes(mux *http.ServeMux) {
	var isDevelopment = os.Getenv("GO_ENV") != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./ui/assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /ui/assets/", http.StripPrefix("/ui/assets/", assetHandler))
}

func main() {
	InitDotEnv()
	mux := http.NewServeMux()
	SetupAssetsRoutes(mux)
	mux.Handle("GET /", templ.Handler(pages.Landing()))
	fmt.Println("Server is running on http://localhost:8090")
	http.ListenAndServe(":8090", mux)

	ctx := context.Background()
	db_url := "user=bth database=testdb"
	pgdb := db.NewDatabase(ctx, db_url)

	// joe, err := db.NewUser(ctx, pgdb, "joe", "j@blow.com")
	joe, err := db.GetUser(ctx, pgdb, "0197ada4-5f8b-77d7-b039-5651eabf19e1")
	if err != nil {
		panic(err)
	}
	udb := db.UserDatabase{&joe, pgdb}
	newprod, err := udb.NewProduct(ctx, "stuff", "3.15")
	if err != nil {
		panic(err)
	}
	fmt.Println(newprod.Name)

	// products, err := pgdb.Query.ListProducts(ctx, joe.ID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(products)
}
