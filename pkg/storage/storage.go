package storage

type Storage interface {
	GetPaste(key string) ([]byte, error)       // Retrieve post data from blob storage
	UploadPaste(key string, data []byte) error // Upload post data to blob storage
	DeletePaste(key string) error
}
