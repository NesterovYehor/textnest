package kafka

type KafkaConfig struct {
	Brokers    []string `yaml:"brokers"`
	Topics     []string `yaml:"topics"`
	GroupID    string   `yaml:"groupID"`
	MaxRetries int      `yaml:"max_retries"`
}

func LoadKafkaConfig(brokers []string, topics []string, groupID string, maxRetries int) *KafkaConfig {
	return &KafkaConfig{
		Brokers:    brokers,
		Topics:     topics,
		GroupID:    groupID,
		MaxRetries: maxRetries,
	}
}
