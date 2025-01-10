package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/services"
)

// ExpiredPasteHandler processes messages for expired pastes
type ExpiredPasteHandler struct {
	pasteService *services.PasteService
}

// NewExpiredPasteHandler initializes a handler for expired paste topic
func NewExpiredPasteHandler(srv *services.PasteService) *ExpiredPasteHandler {
	return &ExpiredPasteHandler{pasteService: srv}
}

// Handle processes the Kafka message to delete expired paste
func (h *ExpiredPasteHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	key := string(msg.Value)
	if err := h.pasteService.DeletePasteByKey(ctx, key); err != nil {
		return err
	}
	return nil
}

