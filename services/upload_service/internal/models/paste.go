package models

type PasteData struct {
	Key     string // Unique identifier for the post (hash)
	Content string // The actual text content of the post
}

type Storage interface {
	GetPaste(key string) (string, error)        // Retrieve post data from blob storage
	UploadPaste(key *PasteData) (string, error) // Upload post data to blob storage
	DeletePaste(key string) error
}
