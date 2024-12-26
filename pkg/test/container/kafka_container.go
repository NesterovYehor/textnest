package container

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

type KafkaSetup struct {
	CleanUp    func()
	BrokerAddr []string
}

func StartKafka(ctx context.Context) (*KafkaSetup, error) {
	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/cp-kafka:latest",
		kafka.WithClusterID("test-cluster"),
	)
	if err != nil {
		return nil, err
	}
	brokerAddr, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		kafkaContainer.Terminate(ctx)
		return nil, err
	}
	log.Printf("Kafka brocker Addr: %v", brokerAddr)
	setUp := KafkaSetup{
		CleanUp: func() {
			kafkaContainer.Terminate(ctx)
		},
		BrokerAddr: brokerAddr,
	}

	log.Printf("Kafka broker is ready at %s", brokerAddr)
	return &setUp, nil
}
