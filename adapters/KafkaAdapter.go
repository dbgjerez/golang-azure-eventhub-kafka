package adapters

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

const (
	saslUser            = "KAFKA_SASL_USER"
	saslPassword        = "KAFKA_SASL_PASSWORD"
	kafkaBrokers        = "KAFKA_BROKERS"
	kafkaConsoumerGroup = "KAFKA_CONSUMER_GROUP"
	kafkaTopic          = "KAFKA_CONSUMER_TOPIC"
)

type KafkaConnection struct {
	kafka sarama.ConsumerGroup
}

func NewConnection() (conn *KafkaConnection) {
	config := KafkaConfig(os.Getenv(saslUser), os.Getenv(saslPassword))
	brokers := strings.Split(os.Getenv(kafkaBrokers), ",")
	consumerGroup := os.Getenv(kafkaConsoumerGroup)
	client, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	conn = &KafkaConnection{client}
	return conn
}

func (consumer *KafkaConnection) Suscribe(ctx context.Context) {
	for {
		err := consumer.kafka.Consume(ctx, []string{os.Getenv(kafkaTopic)}, &Consumer{})
		if err != nil {
			log.Fatalln("Error consuming from group", err)
			os.Exit(1)
		}

		if ctx.Err() != nil {
			//exit for loop
			return
		}
	}
}

func (consumer *KafkaConnection) Close() {
	if err := consumer.kafka.Close(); err != nil {
		log.Fatalln("Error trying to close Kafka client: ", err)
	}
	log.Println("Kafka client closed")
}

func KafkaConfig(user string, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	//config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Net.SASL.Mechanism = "PLAIN"
	config.Net.SASL.Enable = true
	config.Net.SASL.User = user
	config.Net.SASL.Password = password
	//config.Consumer.Group.Session.Timeout = 30 * time.Second
	config.Net.TLS.Enable = true
	return config
}

type Consumer struct{}

func (consumer *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	log.Println("Partition allocation -", s.Claims())
	return nil
}

func (consumer *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	log.Println("Consumer group clean up initiated")
	return nil
}
func (consumer *Consumer) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	for msg := range c.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		log.Println("Message content", string(msg.Value))
		s.MarkMessage(msg, "")
	}
	return nil
}
