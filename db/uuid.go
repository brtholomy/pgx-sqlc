package db

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ReadUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func MakeUUIDv7() (pgtype.UUID, error) {
	googleuuid, err := uuid.NewV7()
	if err != nil {
		return pgtype.UUID{}, err
	}
	return ReadUUID(googleuuid.String())
}

func GetUUIDFromUrlValues(vals url.Values, id string) (pgtype.UUID, error) {
	pidstr := ""
	if vals.Has(id) {
		pidstr = vals.Get(id)
	} else {
		return pgtype.UUID{}, errors.New("id not found in request: " + id)
	}
	pid, err := ReadUUID(pidstr)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pid, nil
}

func GetUUIDFromHeader(header *http.Header, id string) (pgtype.UUID, error) {
	pidstr := ""
	if pidstr = header.Get(id); pidstr == "" {
		return pgtype.UUID{}, errors.New("id not found in header: " + id)
	}
	pid, err := ReadUUID(pidstr)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pid, nil
}
