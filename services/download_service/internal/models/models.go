package models

import "database/sql"

type Models struct {
	Metadata MetadataModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Metadata: MetadataModel{
			DB: db,
		},
	}
}
