package data

import (
	"database/sql"
	"errors"
)

var (
	ErrorRecordNotFound = errors.New("record not found ")
	ErrorEditConflict   = errors.New("edit conflict ")
)

type Models struct {
	Movie MovieModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movie: MovieModel{DB: db},
	}
}
