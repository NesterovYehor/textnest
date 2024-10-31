package models

type PasteData struct {
	Hash    string // Unique identifier for the post (hash)
	Content string // The actual text content of the post
}

type Storage interface {
	GetPaste(hash string) (*PasteData, error)    // Retrieve post data from blob storage
	UploadPaste(data *PasteData) (string, error) // Upload post data to blob storage
	DeletePaste(hash string) error
}
