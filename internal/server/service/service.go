package service

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type Service struct {
	Ctx      context.Context
	Producer sarama.SyncProducer
}

func New(ctx context.Context, producer sarama.SyncProducer) *Service {
	return &Service{
		Ctx:      ctx,
		Producer: producer,
	}
}

func (s *Service) SendMessage(msg *sarama.ProducerMessage) error {
	if _, _, err := s.Producer.SendMessage(msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}
