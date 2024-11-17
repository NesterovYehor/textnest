package config

type KafkaConsumerConfig struct {
	Brokers       []string
	ProducerTopic string
}

func LoadKafkaConfig() *KafkaConsumerConfig {
	return &KafkaConsumerConfig{
		Brokers:       []string{"localhost:9092"},
		ProducerTopic: "expired-paste",
	}
}
