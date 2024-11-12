package storage

type storage interface {
	Uploadpaste(key string, data []byte) error // upload post data to blob storage
}
