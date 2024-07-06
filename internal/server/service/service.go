package service

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Service struct {
	Producer sarama.SyncProducer
}

func (s *Service) SendMessage(msg *sarama.ProducerMessage) error {
	if _, _, err := s.Producer.SendMessage(msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}
