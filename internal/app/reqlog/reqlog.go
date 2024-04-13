package reqlog

import (
	"encoding/json"
	"homework/internal/app/kafka"
)

type Logger struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewLogger(producer *kafka.Producer, consumer *kafka.Consumer) *Logger {
	return &Logger{producer: producer, consumer: consumer}
}

func (l *Logger) Log(msg Message) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return
	}

	l.producer.SendMessage(bytes)
}
