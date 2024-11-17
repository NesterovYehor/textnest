package storage

type Storage interface {
	DeletePaste(key string) error
	DeleteExpiredPastes(keys []string) error
}
