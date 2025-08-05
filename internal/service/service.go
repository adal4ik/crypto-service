package service

import (
	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/pkg/logger"
)

type Service struct {
}

func NewService(repo *repository.Repository, logger logger.Logger) *Service {
	return &Service{
		// SubscriptionService: NewSubscriptionService(repo.SubscriptionRepository, logger),
	}
}
