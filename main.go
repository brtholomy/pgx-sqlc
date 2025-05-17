package main

import (
	"context"
	"fmt"
	"os"

	"pgx-sqlc/db"
	"pgx-sqlc/db/sqlc"
)

func main() {
	ctx := context.Background()
	db_url := "user=bth database=testdb"
	pgdb := db.NewDatabase(ctx, db_url)

	id, err := db.GetUUIDv7()
	newacc, err := pgdb.Query.CreateUser(ctx, sqlc.CreateUserParams{
		ID:    *id,
		Name:  "Bob",
		Email: "bob@j.com",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(newacc.ID)

	acc, err := pgdb.Query.GetUser(ctx, newacc.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(acc.Email)

}
