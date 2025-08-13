package service

import (
	"TestEffectiveMobile/internal/config"
	"TestEffectiveMobile/internal/models"
	"TestEffectiveMobile/internal/repository"
	"TestEffectiveMobile/pkg/logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"regexp"
)

type SubscriptionServiceInterface interface {
	Create(sub *models.Subscription) (string, error)
	Read(id string) (*models.Subscription, error)
	Update(id string, sub *models.UpdateSubscription) error
	Delete(id string) error
	ListSubscriptions(userId string) ([]*models.Subscription, error)
	CalculateSumSubscriptions(userId string, startDate string, endDate string, serviceName string) (int, error)
}

type SubscriptionService struct {
	Repository repository.SubscriptionRepositoryInterface
	ctx        context.Context
	cfg        *config.Config
}

func NewSubscriptionService(repo repository.SubscriptionRepositoryInterface, cfg *config.Config, ctx context.Context) *SubscriptionService {
	return &SubscriptionService{
		Repository: repo,
		cfg:        cfg,
		ctx:        ctx,
	}
}

func (s *SubscriptionService) Create(sub *models.Subscription) (string, error) {
	if sub == nil || sub.ServiceName == "" || sub.Price == 0 || sub.UserId == "" || sub.StartDate == "" || sub.EndDate == "" {
		return "", fmt.Errorf("sub is empty")
	}
	if !IsValidMMYYYY(sub.StartDate) || !IsValidMMYYYY(sub.EndDate) {
		return "", fmt.Errorf("startDate or endDate is not valid")
	}
	sub.Id = uuid.New().String()
	logger.GetLoggerFromCtx(s.ctx).Info(fmt.Sprintf("Create sub: %v", sub))
	return sub.Id, s.Repository.Create(sub)
}

func (s *SubscriptionService) Read(id string) (*models.Subscription, error) {
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	logger.GetLoggerFromCtx(s.ctx).Info(fmt.Sprintf("Read id: %s", id))
	return s.Repository.Read(id)
}

func (s *SubscriptionService) Update(id string, sub *models.UpdateSubscription) error {
	if id == "" || sub == nil || sub.ServiceName == "" || sub.Price == 0 || sub.StartDate == "" || sub.EndDate == "" {
		return fmt.Errorf("sub or id is empty")
	}
	if !IsValidMMYYYY(sub.StartDate) || !IsValidMMYYYY(sub.EndDate) {
		return fmt.Errorf("startDate or endDate is not valid")
	}
	logger.GetLoggerFromCtx(s.ctx).Info(fmt.Sprintf("Update id: %s, sub: %v", id, sub))
	return s.Repository.Update(id, sub)
}

func (s *SubscriptionService) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("id is empty")
	}
	logger.GetLoggerFromCtx(s.ctx).Info(fmt.Sprintf("Delete id: %s", id))
	return s.Repository.Delete(id)
}

func (s *SubscriptionService) ListSubscriptions(userId string) ([]*models.Subscription, error) {
	if userId == "" {
		return nil, fmt.Errorf("user_id is empty")
	}
	logger.GetLoggerFromCtx(s.ctx).Info(fmt.Sprintf("List user_id: %s", userId))
	return s.Repository.ListSubscriptions(userId)
}

func (s *SubscriptionService) CalculateSumSubscriptions(userId string, startDate string, endDate string, serviceName string) (int, error) {
	if startDate == "" || endDate == "" {
		return 0, fmt.Errorf("userId or startDate or endDate or serviceName is empty")
	}
	if !IsValidMMYYYY(startDate) || !IsValidMMYYYY(endDate) {
		return 0, fmt.Errorf("startDate or endDate is not valid")
	}
	logger.GetLoggerFromCtx(s.ctx).Info(fmt.Sprintf("Calculate Sum userId: %s, startDate: %s, endDate: %s, serviceName: %s", userId, startDate, endDate, serviceName))
	return s.Repository.CalculateSumSubscriptions(userId, startDate, endDate, serviceName)
}

func IsValidMMYYYY(dateStr string) bool {
	matched, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-20[2-9][0-9]$`, dateStr)
	return matched
}
