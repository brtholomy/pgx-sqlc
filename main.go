package main

import (
	"context"
	"fmt"

	"pgx-sqlc/db"
	"pgx-sqlc/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

const bob_id string = "0196e091-b831-7209-9827-7f6e41b65296"

func NewUser(ctx context.Context, pgdb *db.Database) (*sqlc.User, error) {
	id, err := db.GetUUIDv7()
	if err != nil {
		return nil, err
	}
	newacc, err := pgdb.Query.CreateUser(ctx, sqlc.CreateUserParams{
		ID:    *id,
		Name:  "Bob",
		Email: "bob@j.com",
	})
	if err != nil {
		return nil, err
	}
	return &newacc, nil
}

func main() {
	ctx := context.Background()
	db_url := "user=bth database=testdb"
	pgdb := db.NewDatabase(ctx, db_url)
	var err error

	// product
	var uid pgtype.UUID
	err = uid.Scan(bob_id)
	if err != nil {
		panic(err)
	}
	pid, err := db.GetUUIDv7()
	if err != nil {
		panic(err)
	}
	var num pgtype.Numeric
	err = num.Scan("1.99")
	if err != nil {
		panic(err)
	}
	newprod, err := pgdb.Query.CreateProduct(ctx, sqlc.CreateProductParams{
		ID:     *pid,
		UserID: uid,
		Name:   "gum",
		Price:  num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(newprod.Name)

	products, err := pgdb.Query.ListProducts(ctx, uid)
	if err != nil {
		panic(err)
	}
	fmt.Println(products)
}
