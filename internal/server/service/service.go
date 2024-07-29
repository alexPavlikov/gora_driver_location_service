package service

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/gora_driver_location_service/internal/models"
	"github.com/alexPavlikov/gora_driver_location_service/internal/server/repository"
	"github.com/vmihailenco/msgpack/v5"
)

type Service struct {
	Repo *repository.Repo
}

func New(repo *repository.Repo) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) StoreMessage(ctx context.Context, cord models.Cord) error {

	var id = ctx.Value("Driver_id")

	value, err := msgpack.Marshal(&cord)
	if err != nil {
		return fmt.Errorf("msgpack marshal err: %w", err)
	}

	var msg = models.CoordinatesPayload{
		Key:   id.(int),
		Value: value,
	}

	if err := s.Repo.SendMessage(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (s *Service) ReadMessage(ctx context.Context) ([]models.Cord, error) {
	cords, err := s.Repo.ReadMessageFromKafka()
	if err != nil {
		return nil, err
	}

	return cords, nil
}
