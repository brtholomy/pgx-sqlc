package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"pgx-sqlc/db"
	"pgx-sqlc/ui/assets"
	"pgx-sqlc/ui/pages"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
)

const RENDER_ENV_PATH = "/etc/secrets/.env"

// FIXME: yeah
const LOCALJOE = "0197ada4-5f8b-77d7-b039-5651eabf19e1"
const REMOTEJOE = "0197d266-380a-741a-85f8-04fe8358be6b"

func initDotEnv() error {
	var err error
	if err = godotenv.Load(RENDER_ENV_PATH); err == nil {
		return nil
	} else if err = godotenv.Load(); err == nil {
		return nil
	}
	return errors.New("Could not load .env file")
}

func setupAssetsRoutes(mux *http.ServeMux, isdev bool) {
	asset_handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isdev {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isdev {
			fs = http.FileServer(http.Dir("./ui/assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /ui/assets/", http.StripPrefix("/ui/assets/", asset_handler))
}

func main() {
	log.SetFlags(log.Lshortfile)
	if err := initDotEnv(); err != nil {
		panic(err)
	}
	var isdev = os.Getenv("GO_ENV") != "production"

	////////////////////////////////////////
	// setup DB
	ctx := context.Background()
	db_url := "user=bth database=testdb"
	if !isdev {
		if tmp, ok := os.LookupEnv("RENDER_INTERNAL_DB"); ok {
			db_url = tmp
		} else {
			log.Fatal("failed to load RENDER_INTERNAL_DB")
		}
	}
	log.Printf("db_url: %#v\n", db_url)
	pgdb, err := db.NewDatabase(ctx, db_url)
	if err != nil {
		panic(err)
	}

	// TODO: user login flow
	// joe, err := db.NewUser(ctx, pgdb, "joe", "j@blow.com")
	UUID := LOCALJOE
	if !isdev {
		UUID = REMOTEJOE
	}
	joe, err := db.GetUser(ctx, pgdb, UUID)
	if err != nil {
		panic(err)
	}
	udb := db.UserDatabase{&joe, pgdb}

	getProductsH, err := db.InitHandler(&udb, db.GetProducts)
	if err != nil {
		panic(err)
	}
	postProductsH, err := db.InitHandler(&udb, db.PostProducts)
	if err != nil {
		panic(err)
	}

	// ////////////////////////////////////////
	// // QBO setup
	// c, err := qbo.SetupQboClient()
	// if err != nil {
	// 	panic(err)
	// }
	// geth, err := qbo.InitHandler(c, qbo.GetInvoice)
	// if err != nil {
	// 	panic(err)
	// }
	// posth, err := qbo.InitHandler(c, qbo.PostInvoice)
	// if err != nil {
	// 	panic(err)
	// }

	////////////////////////////////////////
	// handlers
	mux := http.NewServeMux()
	setupAssetsRoutes(mux, isdev)

	// default
	mux.Handle("GET /", templ.Handler(pages.Landing()))
	// DB
	mux.Handle("GET /products", getProductsH)
	mux.Handle("POST /products", postProductsH)
	// QBO
	// mux.Handle("GET /qbo", geth)
	// mux.Handle("POST /qbo", posth)

	// ListenAndServe
	if isdev {
		fmt.Println("Server is running on http://localhost:8090")
	}
	http.ListenAndServe(":8090", mux)

}
