package kafka

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/config"
)

type KafkaProducer struct {
	asyncProducer sarama.AsyncProducer
	kafkaConfig   config.KafkaConfig
}

func NewProducer(kafkaConfig config.KafkaConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Retry.Max = 5

	producer, err := sarama.NewAsyncProducer(kafkaConfig.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		asyncProducer: producer,
		kafkaConfig:   kafkaConfig,
	}, nil
}

func SendExpiredKeysToKafka(producer *KafkaProducer, message string) error {
	return producer.produceMessages(message)
}

func (producer *KafkaProducer) produceMessages(messageValue string) error {
	message := &sarama.ProducerMessage{
		Topic: producer.kafkaConfig.ProducerTopic,
		Value: sarama.StringEncoder(messageValue),
	}

	// Asynchronously produce the message
	select {
	case producer.asyncProducer.Input() <- message:
		log.Println("New Message produced to topic", producer.kafkaConfig.ProducerTopic)
	case err := <-producer.asyncProducer.Errors():
		log.Println("Error producing message:", err)
		return err
	}

	return nil
}


func (producer *KafkaProducer) Close() error {
	// Close the producer to drain messages and release resources
	return producer.asyncProducer.Close()
}
