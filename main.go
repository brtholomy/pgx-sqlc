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

	products, err := pgdb.Query.ListProducts(ctx, joe.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(products)
}
