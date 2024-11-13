package storage

type Storage interface {
	DownloadPaste(key string) ([]byte, error) // upload post data to blob storage
}
