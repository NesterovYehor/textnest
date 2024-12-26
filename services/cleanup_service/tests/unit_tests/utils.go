package testutils

import (
	"database/sql"

	"github.com/stretchr/testify/assert"
)

// VerifyRowExists checks if a row with the given key exists in the metadata table.
func VerifyRowExists(t assert.TestingT, db *sql.DB, key string) bool {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM metadata WHERE key = $1`, key).Scan(&count)
	assert.NoError(t, err)
	return count > 0
}
