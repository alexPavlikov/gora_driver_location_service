package kafka

import (
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

// GetProducer Функция получения producer
func GetProducer(address ...string) (producer sarama.SyncProducer, err error) {
	producer, err = sarama.NewSyncProducer(address, nil) //....
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %w", err)
	}

	return producer, nil
}

// WriteMessage Функция записи сообщения в kafka
func WriteMessage(producer sarama.SyncProducer, msg *sarama.ProducerMessage) error {
	_, _, err := producer.SendMessage(msg)
	if err != nil {
		slog.Error("write message to kafka error - " + err.Error())
		return err
	}

	return nil
}
