package storage

type Storage interface {
	UploadPaste(key string, data []byte) error // Upload post data to blob storage
}
