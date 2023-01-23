package dbexplorer

import (
	"hw6coursera/dto"
	"hw6coursera/repository"
	"log"
)

type SchemeParserExplorer struct {
	repoExplorer repository.Explorer
}

// ParseSchema implements SchemeParser
func (s *SchemeParserExplorer) ParseSchema() (dto.Schema, error) {
	log.Println("getting tables")
	tableNames, err := s.repoExplorer.GetTableNames()
	if err != nil {
		return nil, err
	}

	sch := make(dto.Schema, len(tableNames))

	for _, tableName := range tableNames {
		t := dto.Table{}
		log.Printf("parsing colunms in table: %s", tableName)
		cols, err := s.repoExplorer.GetColumns(tableName)
		if err != nil {
			return nil, err
		}
		t.Name = tableName
		t.Columns = cols
		sch[tableName] = t
	}
	return sch, nil
}

func newSchemeParser(r *repository.Repository) *SchemeParserExplorer {
	return &SchemeParserExplorer{
		repoExplorer: r.Explorer,
	}
}
