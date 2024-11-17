package kafka

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	asyncProducer sarama.AsyncProducer
	kafkaConfig   KafkaConfig
	ctx           context.Context
}

func NewProducer(kafkaConfig KafkaConfig, ctx context.Context) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Retry.Max = kafkaConfig.MaxRetries

	producer, err := sarama.NewAsyncProducer(kafkaConfig.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		asyncProducer: producer,
		kafkaConfig:   kafkaConfig,
		ctx:           ctx,
	}, nil
}

func (producer *KafkaProducer) ProduceMessages(messageValue string, topic string) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageValue),
	}

	for attempts := 0; attempts < producer.kafkaConfig.MaxRetries; attempts++ {
		select {
		case producer.asyncProducer.Input() <- message:
			log.Println("New Message produced to topic", topic)
			return nil
		case err := <-producer.asyncProducer.Errors():
			if attempts < producer.kafkaConfig.MaxRetries-1 {
                time.Sleep(time.Duration(math.Pow(2, float64(attempts))) * time.Second)
				continue
			}
			log.Println("Error producing message:", err)
			return err
		}
	}

	return nil
}

func (producer *KafkaProducer) Close() error {
	// Close the producer to drain messages and release resources
	return producer.asyncProducer.Close()
}
