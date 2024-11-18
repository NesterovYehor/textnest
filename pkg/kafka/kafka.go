package kafka

type KafkaConfig struct {
	Brokers    []string
	Topics     map[string]string
	GroupID    string
	MaxRetries int
}

func LoadKafkaConfig(brokers []string, topics map[string]string, groupID string, maxRetries int) *KafkaConfig {
	return &KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		GroupID:    groupID,
		MaxRetries: maxRetries,
	}
}
