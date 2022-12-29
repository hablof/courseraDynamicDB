package dbexplorer

import (
	"hw6coursera/internal"
	"hw6coursera/repository"
)

type SchemeParser interface {
	ParseSchema() (internal.Schema, error)
}

type DBexplorer struct {
	SchemeParser
}

func NewDbExplorer(r *repository.Repository) *DBexplorer {
	return &DBexplorer{
		SchemeParser: newSchemeParser(r),
	}
}
