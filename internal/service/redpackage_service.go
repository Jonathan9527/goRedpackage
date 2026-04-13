package service

import (
	"context"
	"errors"
	"math/rand"

	"learnGO/internal/database"
	"learnGO/internal/repository"

	"github.com/shopspring/decimal"
)

type RedPackageService struct {
	userRepository *repository.UserRepository
	publisher      *database.RabbitMQPublisher
}

type RedPackageCreatedMessage struct {
	Account        string    `json:"account"`
	TotalAmount    string    `json:"total_amount"`
	RedPackageList []float64 `json:"red_package_list"`
}

func NewRedPackageService(userRepository *repository.UserRepository, publisher *database.RabbitMQPublisher) *RedPackageService {
	return &RedPackageService{
		userRepository: userRepository,
		publisher:      publisher,
	}
}

func (s *RedPackageService) CreateRedPackage(ctx context.Context, account string, redAmount decimal.Decimal) ([]float64, error) {
	user, err := s.userRepository.FindByAccount(ctx, account)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Balance.LessThan(redAmount) {
		return nil, errors.New("insufficient balance")
	}

	if err := s.userRepository.UpdateBalance(ctx, user, redAmount); err != nil {
		return nil, err
	}

	redPackageList := makeRedPackageList(redAmount.IntPart(), 5)
	if err := s.publisher.PublishJSON("red_package.created", RedPackageCreatedMessage{
		Account:        account,
		TotalAmount:    redAmount.StringFixed(2),
		RedPackageList: redPackageList,
	}); err != nil {
		return nil, err
	}

	return redPackageList, nil
}

func makeRedPackageList(totalAmount int64, totalNum int) []float64 {
	result := make([]float64, totalNum)

	remainAmount := totalAmount * 100 // 转为分
	remainNum := totalNum

	for i := 0; i < totalNum; i++ {
		if remainNum == 1 {
			result[i] = float64(remainAmount) / 100 // 转回元
		} else {
			maxAmount := remainAmount / int64(remainNum) * 2
			if maxAmount <= 0 {
				maxAmount = 1
			}
			randomAmount := rand.Int63n(maxAmount) + 1
			result[i] = float64(randomAmount) / 100 // 转回元
			remainAmount -= randomAmount
		}
		remainNum--
	}

	return result
}
