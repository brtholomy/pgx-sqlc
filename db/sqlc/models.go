// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Invoice struct {
	ID            pgtype.UUID        `db:"id" json:"id"`
	UserID        pgtype.UUID        `db:"user_id" json:"user_id"`
	InvoiceNumber int32              `db:"invoice_number" json:"invoice_number"`
	Total         pgtype.Numeric     `db:"total" json:"total"`
	Created       pgtype.Timestamptz `db:"created" json:"created"`
	Modified      pgtype.Timestamptz `db:"modified" json:"modified"`
}

type InvoiceItem struct {
	ID        pgtype.UUID    `db:"id" json:"id"`
	UserID    pgtype.UUID    `db:"user_id" json:"user_id"`
	ProductID pgtype.UUID    `db:"product_id" json:"product_id"`
	InvoiceID pgtype.UUID    `db:"invoice_id" json:"invoice_id"`
	Amount    pgtype.Numeric `db:"amount" json:"amount"`
}

type Product struct {
	ID     pgtype.UUID    `db:"id" json:"id"`
	UserID pgtype.UUID    `db:"user_id" json:"user_id"`
	Name   string         `db:"name" json:"name"`
	Price  pgtype.Numeric `db:"price" json:"price"`
}

type User struct {
	ID    pgtype.UUID `db:"id" json:"id"`
	Email string      `db:"email" json:"email"`
	Name  string      `db:"name" json:"name"`
}
