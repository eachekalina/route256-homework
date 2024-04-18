package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"golang.org/x/sync/errgroup"
	"time"
)

type MessageHandler func(msg []byte)

type Consumer struct {
	client         sarama.Consumer
	topic          string
	messageHandler MessageHandler
	ready          chan bool
}

func NewConsumer(brokers []string, topic string, messageHandler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 10 * time.Second
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	client, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:         client,
		topic:          topic,
		messageHandler: messageHandler,
		ready:          make(chan bool),
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	patritions, err := c.client.Partitions(c.topic)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)

	for _, p := range patritions {
		pc, err := c.client.ConsumePartition(c.topic, p, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case msg, ok := <-pc.Messages():
					if !ok {
						return nil
					}
					c.messageHandler(msg.Value)
				}
			}
		})
	}
	close(c.ready)

	return eg.Wait()
}

func (c *Consumer) Ready() <-chan bool {
	return c.ready
}
