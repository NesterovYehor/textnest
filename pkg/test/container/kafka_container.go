package container

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

type KafkaContainerOpts struct {
	ClusterID         string           // Kafka Cluster ID (default: "test-cluster")
	BrokerPort        int              // Port for Kafka broker
	Topics            map[string]int32 // Topics to create with their partitions
	ReplicationFactor int16            // Replication factor for topics
}

func StartKafka(ctx context.Context, opts *KafkaContainerOpts) (*kafka.KafkaContainer, string, error) {
	// Set default options if not provided
	if opts.ClusterID == "" {
		opts.ClusterID = "test-cluster"
	}
	if opts.BrokerPort == 0 {
		opts.BrokerPort = 9092
	}
	if opts.ReplicationFactor == 0 {
		opts.ReplicationFactor = 1
	}

	// Start the Kafka container
	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/cp-kafka:latest",
		kafka.WithClusterID(opts.ClusterID),
		testcontainers.WithHostPortAccess(opts.BrokerPort),
	)
	if err != nil {
		return nil, "", err
	}

	// Retrieve container's host and port
	host, err := kafkaContainer.Host(ctx)
	if err != nil {
		kafkaContainer.Terminate(ctx)
		return nil, "", err
	}
	mappedPort, err := kafkaContainer.MappedPort(ctx, nat.Port(opts.BrokerPort))
	if err != nil {
		kafkaContainer.Terminate(ctx)
		return nil, "", err
	}

	// Construct the broker address (host:port)
	brokerAddr := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	// If topics are specified, create them
	if len(opts.Topics) > 0 {
		err = createTopics(brokerAddr, opts.Topics, opts.ReplicationFactor)
		if err != nil {
			kafkaContainer.Terminate(ctx)
			return nil, "", err
		}
	}

	return kafkaContainer, brokerAddr, nil
}

func createTopics(brokerAddr string, topics map[string]int32, replicationFactor int16) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{brokerAddr}, config)
	if err != nil {
		return err
	}
	defer admin.Close()

	for topic, partitions := range topics {
		topicDetail := &sarama.TopicDetail{
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		}
		err = admin.CreateTopic(topic, topicDetail, false)
		if err != nil {
			log.Printf("Failed to create topic %s: %v", topic, err)
			return err
		}
		log.Printf("Topic %s created successfully", topic)
	}

	return nil
}
