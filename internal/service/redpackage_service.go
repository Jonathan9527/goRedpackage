package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"learnGO/internal/database"
	"learnGO/internal/repository"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type RedPackageService struct {
	userRepository *repository.UserRepository
	publisher      *database.RabbitMQPublisher
	redisClient    *redis.Client
}

type RedPackageCreatedMessage struct {
	Account        string    `json:"account"`
	TotalAmount    string    `json:"total_amount"`
	RedPackageList []float64 `json:"red_package_list"`
}

type CachedRedPackage struct {
	Account        string    `json:"account"`
	RedPackageList []float64 `json:"red_package_list"`
}

func NewRedPackageService(userRepository *repository.UserRepository, publisher *database.RabbitMQPublisher, redisClient *redis.Client) *RedPackageService {
	return &RedPackageService{
		userRepository: userRepository,
		publisher:      publisher,
		redisClient:    redisClient,
	}
}

func (s *RedPackageService) CreateRedPackage(ctx context.Context, account string, redAmount decimal.Decimal, number int) ([]float64, error) {
	user, err := s.userRepository.FindByAccount(ctx, account)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Balance.LessThan(redAmount) {
		return nil, errors.New("insufficient balance")
	}

	pkgId, err := s.userRepository.UpdateBalance(ctx, user, redAmount)
	if err != nil {
		return nil, err
	}
	// if err := s.userRepository.UpdateBalance(ctx, user, redAmount); err != nil {
	// 	return nil, err
	// }
	redPackageList := makeRedPackageList(redAmount.IntPart(), number)
	if err := s.cacheRedPackageList(ctx, account, redPackageList, pkgId); err != nil {
		return nil, err
	}
	if err := s.publisher.PublishJSON("red_package.created", RedPackageCreatedMessage{
		Account:        account,
		TotalAmount:    redAmount.StringFixed(2),
		RedPackageList: redPackageList,
	}); err != nil {
		return nil, err
	}

	return redPackageList, nil
}

func (s *RedPackageService) cacheRedPackageList(ctx context.Context, account string, redPackageList []float64, pkgId int64) error {
	if s.redisClient == nil {
		return nil
	}

	// // body, err := json.Marshal(gin.H{
	// body, err := json.Marshal(CachedRedPackage{
	// 	Account:        account,
	// 	RedPackageList: redPackageList,
	// })
	// if err != nil {
	// 	return err
	// }
	//写入总数key
	totalkey := fmt.Sprintf("redpackage_total:%d", pkgId)
	if err := s.redisClient.Set(ctx, totalkey, len(redPackageList), 24*7*time.Hour).Err(); err != nil {

	}
	redPackageListJSON, _ := json.Marshal(redPackageList)

	//写入红包详情key
	detailKey := fmt.Sprintf("redpackage_detail:%d", pkgId)
	if err := s.redisClient.HSet(ctx, detailKey, map[string]interface{}{
		"account": account,
		"list":    redPackageListJSON,
		"total":   len(redPackageList),
	}).Err(); err != nil {
		return err
	}
	//写入红包列表key
	redListkey := fmt.Sprintf("redpackage_list:%d", pkgId)
	values := make([]interface{}, len(redPackageList))
	for i, v := range redPackageList {
		values[i] = v
	}
	return s.redisClient.LPush(ctx, redListkey, values...).Err()
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
