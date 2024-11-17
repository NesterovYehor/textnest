package kafka

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/services/key_generation_service/internal/repository"
)

type KeyReallocatorConsumer struct {
	consumer sarama.Consumer
	repo     *repository.KeymanagerRepo
	topic    string
}

// NewKeyReallocatorConsumer initializes a new KeyReallocatorConsumer.
func NewKeyReallocatorConsumer(brokers []string, topic string, repo *repository.KeymanagerRepo) (*KeyReallocatorConsumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KeyReallocatorConsumer{
		consumer: consumer,
		repo:     repo,
		topic:    topic,
	}, nil
}

func (c *KeyReallocatorConsumer) Start(ctx context.Context) error {
	var partitionConsumer sarama.PartitionConsumer
	var err error

	// Retry loop for restarting the consumer
	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer shutting down...")
			if partitionConsumer != nil {
				partitionConsumer.Close()
			}
			return nil

		default:
			// Attempt to consume partition
			partitionConsumer, err = c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Error consuming partition: %v. Retrying...", err)
				time.Sleep(5 * time.Second) // Backoff before retrying
				continue
			}
			log.Println("Consumer started successfully.")

			// Process messages
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					log.Printf("Message received: %s\n", string(msg.Value))
					if err := c.repo.ReallocateKey(string(msg.Value)); err != nil {
						log.Printf("Error reallocating key: %v\n", err)
					}
				case <-ctx.Done():
					log.Println("Stopping partition consumer...")
					partitionConsumer.Close()
					return nil
				}
			}
		}
	}
}
