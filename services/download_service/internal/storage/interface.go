package storage

type storage interface {
	DownloadPaste(key string) ([]byte, error) // upload post data to blob storage
}
