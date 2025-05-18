package db

import (
	"context"
	"pgx-sqlc/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewUser(ctx context.Context, pgdb *Database, name, email string) (*sqlc.User, error) {
	id, err := GetUUIDv7()
	if err != nil {
		return nil, err
	}
	newacc, err := pgdb.Query.CreateUser(ctx, sqlc.CreateUserParams{
		ID:    *id,
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}
	return &newacc, nil
}

// Another wrapper.
// TODO: Not convinced this is the right abstraction. But surely I'll be doing most things on the
// behalf of a single user.
type UserDatabase struct {
	User *sqlc.User
	DB   *Database
}

func (udb *UserDatabase) NewProduct(ctx context.Context, name, price string) (*sqlc.Product, error) {
	pid, err := GetUUIDv7()
	if err != nil {
		return nil, err
	}
	var num pgtype.Numeric
	err = num.Scan(price)
	if err != nil {
		return nil, err
	}
	newprod, err := udb.DB.Query.CreateProduct(ctx, sqlc.CreateProductParams{
		ID:     *pid,
		UserID: udb.User.ID,
		Name:   name,
		Price:  num,
	})
	if err != nil {
		return nil, err
	}
	return &newprod, nil
}
