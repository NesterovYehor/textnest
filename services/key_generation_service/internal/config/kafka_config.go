package config

type KafkaConfig struct {
	Brokers       []string
	ConsumerTopic string
	GroupID       string
}

func LoadKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers:       []string{"localhost:9092"},
		ConsumerTopic: "",
		GroupID:       "expiration-service",
	}
}
