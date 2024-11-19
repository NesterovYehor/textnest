package handlers

import (
	"context"

	"github.com/IBM/sarama"
	jsonlog "github.com/NesterovYehor/TextNest/pkg/logger"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/services"
)

func HandleDeleteExpiredPaste(ctx context.Context, message *sarama.ConsumerMessage, srv *services.PasteService, log *jsonlog.Logger, bucketName string) error {
	key := string(message.Value)

	if err := srv.DeletePasteByKey(ctx, key, bucketName); err != nil {
		log.PrintError(ctx, err, nil)
		return err
	}

	log.PrintInfo(ctx, "Expired paste is deleted", nil)

	return nil
}
