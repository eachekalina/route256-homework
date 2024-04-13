package kafka

import (
	"github.com/IBM/sarama"
	"homework/internal/app/logger"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	log      logger.Logger
}

func NewProducer(brokers []string, log logger.Logger, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Net.MaxOpenRequests = 1
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, topic: topic, log: log}, nil
}

func (p *Producer) SendMessage(msg []byte) {
	prodMsg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
	}
	p.log.Log("Sending message...")
	_, _, err := p.producer.SendMessage(prodMsg)
	if err != nil {
		p.log.Log("%v", err)
	} else {
		p.log.Log("Sent message!")
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
