package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"learnGO/internal/database"
	"learnGO/internal/service"
)

type RedPackageConsumer struct {
	publisher         *database.RabbitMQPublisher
	redPackageService *service.RedPackageService
}

func NewRedPackageConsumer(publisher *database.RabbitMQPublisher, redPackageService *service.RedPackageService) *RedPackageConsumer {
	return &RedPackageConsumer{
		publisher:         publisher,
		redPackageService: redPackageService,
	}
}

func (c *RedPackageConsumer) Start(ctx context.Context) error {

	startConsumersErr := c.startConsumers(c.publisher, 5)
	if startConsumersErr != nil {
		return startConsumersErr
	}

	// deliveries, err := c.publisher.Consume("red_package_pool.created")
	// if err != nil {
	// 	return err
	// }
	// log.Printf("red package consumer started, waiting for messages...")
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case msg, ok := <-deliveries:
	// 			if !ok {
	// 				return
	// 			}

	// 			var payload service.RedPackageCreatedMessage
	// 			if err := json.Unmarshal(msg.Body, &payload); err != nil {
	// 				log.Printf("red package consumer unmarshal failed: %v", err)
	// 				continue
	// 			}
	// 			c.dealRedPackageCreatedMessage(payload)
	// 			msg.Ack(false)
	// 			// log.Printf("red package consumer received: account=%s total=%s id=%s time=%s", payload.Account, payload.Amount, payload.Red_package_id, payload.Create_time)
	// 		}
	// 	}
	// }()

	return nil
}

func (c *RedPackageConsumer) startConsumers(publisher *database.RabbitMQPublisher, workerCount int) error {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			deliveries, err := publisher.Consume("red_package_pool.created")
			if err != nil {
				log.Printf("consumer %d start failed: %v", workerID, err)
				return
			}

			for msg := range deliveries {
				log.Printf("consumer %d received: %s", workerID, string(msg.Body))
				var payload service.RedPackageCreatedMessage
				if err := json.Unmarshal(msg.Body, &payload); err != nil {
					log.Printf("consumer %d unmarshal failed: %v", workerID, err)
					// msg.Nack(false, false) // 处理失败，拒绝消息且不重回队列
					continue
				}
				// 处理业务逻辑
				c.dealRedPackageCreatedMessage(payload)
				// 成功则 ack
				msg.Ack(false)
				// 失败则 nack/requeue
			}
		}(i + 1)
	}

	return nil
}

func (c *RedPackageConsumer) StopstartConsumersWithWorkPool(publisher *database.RabbitMQPublisher, workerCount int) {

}

func (c *RedPackageConsumer) dealRedPackageCreatedMessage(msg service.RedPackageCreatedMessage) {
	time.Sleep(2 * time.Second) // 模拟处理时间
	log.Printf("red package consumer received: account=%s total=%s id=%s time=%s", msg.Account, msg.Amount, msg.Red_package_id, msg.Create_time)
	c.redPackageService.DealUserRedPackageCreatedMessage(context.Background(), msg)
	// c.redPackageService.dealUserRedPackageCreatedMessage(context.Background(), msg)
	// 这里可以添加更多的业务逻辑，比如记录日志、更新统计数据等

}
