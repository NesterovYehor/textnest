package storage

type Storage interface {
	UploadPaste(key string, data []byte) error // upload post data to blob storage
}
