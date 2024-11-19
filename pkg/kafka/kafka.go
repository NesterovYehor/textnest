package kafka

type KafkaConfig struct {
	Brokers    []string
	Topics     []string
	GroupID    string
	MaxRetries int
}

func LoadKafkaConfig(brokers []string, topics []string, groupID string, maxRetries int) *KafkaConfig {
	return &KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		GroupID:    groupID,
		MaxRetries: maxRetries,
	}
}
