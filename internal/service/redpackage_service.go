package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
	Account        string `json:"account"`
	Amount         string `json:"amount"`
	Create_time    string `json:"create_time"`
	Red_package_id string `json:"red_package_id"`
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
	redPackageList := makeRedPackageList(redAmount.IntPart(), number)
	if err := s.cacheRedPackageList(ctx, account, redPackageList, pkgId); err != nil {
		return nil, err
	}

	return redPackageList, nil
}

func (s *RedPackageService) cacheRedPackageList(ctx context.Context, account string, redPackageList []float64, pkgId int64) error {
	if s.redisClient == nil {
		return nil
	}
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

func (s *RedPackageService) DealUserRedPackageCreatedMessage(ctx context.Context, msg RedPackageCreatedMessage) {

	log.Printf("处理红包发放消息: account=%s amount=%s id=%s time=%s", msg.Account, msg.Amount, msg.Red_package_id, msg.Create_time)
	// 这里可以添加一些处理逻辑，比如记录日志、更新统计数据等
	uid, err := strconv.ParseInt(msg.Account, 10, 64)
	if err != nil {
		log.Printf("处理红包发放消息失败UID解析: account=%s, err=%v", msg.Account, err)
		return
	}
	user, err := s.userRepository.FindByUid(ctx, uid)
	if err != nil {
		log.Printf("处理红包发放消息失败: account=%s, err=%v", msg.Account, err)
		return
	}
	redAmount, _ := decimal.NewFromString(msg.Amount)
	pkgId, err := s.userRepository.UpdateGetUserBalance(ctx, user, redAmount, msg.Create_time)
	if err != nil {
		log.Printf("处理红包发放消息失败: account=%s, err=%v", msg.Account, err)
		return
	}
	log.Printf("处理红包发放消息成功: account=%s, pkgId=%d", msg.Account, pkgId)
}
