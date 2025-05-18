package main

import (
	"context"
	"fmt"

	"pgx-sqlc/db"
)

func main() {
	ctx := context.Background()
	db_url := "user=bth database=testdb"
	pgdb := db.NewDatabase(ctx, db_url)
	var err error

	newuser, err := db.NewUser(ctx, pgdb, "joe", "j@blow.com")
	udb := db.UserDatabase{newuser, pgdb}
	newprod, err := udb.NewProduct(ctx, "crap", "0.15")
	fmt.Println(newprod.Name)

	products, err := pgdb.Query.ListProducts(ctx, newuser.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(products)
}
