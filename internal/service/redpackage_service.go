package service

import (
	"context"
	"errors"
	"math/rand"

	"learnGO/internal/repository"

	"github.com/shopspring/decimal"
)

type RedPackageService struct {
	userRepository *repository.UserRepository
}

func NewRedPackageService(userRepository *repository.UserRepository) *RedPackageService {
	return &RedPackageService{
		userRepository: userRepository,
	}
}

func (s *RedPackageService) CreateRedPackage(ctx context.Context, account string, redAmount decimal.Decimal) ([]float64, error) {
	user, err := s.userRepository.FindByAccount(ctx, account)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if decimal.RequireFromString(user.Balance).LessThan(redAmount) {
		return nil, errors.New("insufficient balance")
	}

	// if err := s.userRepository.UpdateBalance(ctx, account, decimal.RequireFromString(user.Balance).Sub(redAmount)); err != nil {
	// 	return nil, err
	// }

	redPackageList := makeRedPackageList(redAmount.IntPart(), 5)
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
