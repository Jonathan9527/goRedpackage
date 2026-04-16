package service

import (
	"context"
	"errors"
	"fmt"
	"learnGO/internal/database"
	"time"

	"github.com/redis/go-redis/v9"
)

type GetRedPackageService struct {
	redisClient *redis.Client
	publisher   *database.RabbitMQPublisher
}

type GetRedPackageResult struct {
	Amount string `json:"amount"`
	Status string `json:"status"`
}

func NewGetRedPackageService(redisClient *redis.Client, publisher *database.RabbitMQPublisher) *GetRedPackageService {
	return &GetRedPackageService{
		redisClient: redisClient,
		publisher:   publisher,
	}
}

type MqRedPackageCreatedMessage struct {
	Account      string `json:"account"`
	Amount       string `json:"amount"`
	CreateTime   string `json:"create_time"`
	RedPackageId string `json:"red_package_id"`
}

func (s *GetRedPackageService) GetRedPackage(ctx context.Context, redPackageID string, userID string) (*GetRedPackageResult, error) {
	if s.redisClient == nil {
		return nil, errors.New("redis client is nil")
	}

	totalKey := fmt.Sprintf("redpackage_total:%s", redPackageID)
	listKey := fmt.Sprintf("redpackage_list:%s", redPackageID)
	poolKey := fmt.Sprintf("redpackage_user:%s", redPackageID)

	luaScript := redis.NewScript(`
local isget = redis.call("ZSCORE", KEYS[3], ARGV[1])
if isget then
	return {0, "您已经领取" .. isget .. "元"}
end
local total = tonumber(redis.call("GET", KEYS[1]) or "0")
if total <= 0 then
    return {0, "红包发放完成1"}
end
local amount = redis.call("RPOP", KEYS[2])
if not amount then
    return {0, "红包发放完成2"}
end
redis.call("DECR", KEYS[1])
redis.call("ZADD", KEYS[3],tonumber(amount), ARGV[1])
redis.call("EXPIRE", KEYS[3], tonumber(ARGV[2]))
return {amount, "ok"}
`)

	result, err := luaScript.Run(ctx, s.redisClient, []string{totalKey, listKey, poolKey}, userID, 86400).Result()
	if err != nil {
		return nil, err
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return nil, fmt.Errorf("unexpected lua result: %#v", result)
	}

	amount := fmt.Sprint(values[0])
	status := fmt.Sprint(values[1])
	if status != "ok" {
		return &GetRedPackageResult{
			Amount: amount,
			Status: status,
		}, nil
	}
	if err := s.publisher.PublishJSON("red_package.created", MqRedPackageCreatedMessage{
		Account:      userID,
		Amount:       amount,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		RedPackageId: redPackageID,
	}); err != nil {
		return nil, err
	}
	return &GetRedPackageResult{
		Amount: amount,
		Status: status,
	}, nil
}
