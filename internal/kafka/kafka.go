package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// Функция получения producer
func GetProducer(address ...string) (producer sarama.SyncProducer, err error) {
	producer, err = sarama.NewSyncProducer(address, nil) //....
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %w", err)
	}

	return producer, nil
}
