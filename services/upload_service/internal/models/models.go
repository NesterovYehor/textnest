package models

import "database/sql"

type Models struct {
	MetaData MetadataModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		MetaData: MetadataModel{
			db: db,
		},
	}
}
