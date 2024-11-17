package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	brokers  []string
	groupID  string
	topics   []string
	handlers map[string]MessageHandler
	ctx      context.Context
	consumer sarama.ConsumerGroup
}

// MessageHandler defines the function signature for message handlers.
type MessageHandler func(*sarama.ConsumerMessage) error

// NewKafkaConsumer initializes a Kafka consumer.
func NewKafkaConsumer(brokers []string, groupID string, topics []string, handlers map[string]MessageHandler, ctx context.Context) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		brokers:  brokers,
		groupID:  groupID,
		topics:   topics,
		handlers: handlers,
		ctx:      ctx,
		consumer: consumerGroup,
	}, nil
}

// Start begins consuming messages and dispatches them to appropriate handlers.
func (kc *KafkaConsumer) Start() error {
	handler := &consumerGroupHandler{handlers: kc.handlers}
	for {
		select {
		case <-kc.ctx.Done():
			log.Println("Kafka consumer shutting down...")
			return kc.consumer.Close()
		default:
			if err := kc.consumer.Consume(kc.ctx, kc.topics, handler); err != nil {
				log.Printf("Error in consumer group: %v", err)
			}
		}
	}
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler.
type consumerGroupHandler struct {
	handlers map[string]MessageHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if handler, exists := h.handlers[message.Topic]; exists {
			if err := handler(message); err != nil {
				log.Printf("Error processing message from topic %s: %v", message.Topic, err)
			}
			session.MarkMessage(message, "")
		} else {
			log.Printf("No handler found for topic: %s", message.Topic)
		}
	}
	return nil
}
