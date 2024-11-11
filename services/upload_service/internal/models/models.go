package models

import "database/sql"

type Models struct {
	Paste PasteModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Paste: PasteModel{
			db: db,
		},
	}
}
