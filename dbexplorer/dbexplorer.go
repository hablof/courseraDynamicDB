package dbexplorer

import (
	"hw6coursera/dto"
	"hw6coursera/repository"
)

type SchemeParser interface {
	ParseSchema() (dto.Schema, error)
}

type DBexplorer struct {
	SchemeParser
}

func NewDbExplorer(r *repository.Repository) *DBexplorer {
	return &DBexplorer{
		SchemeParser: newSchemeParser(r),
	}
}
