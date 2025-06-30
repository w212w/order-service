package kafka

import (
	"context"
	"encoding/json"
	"log"
	"order-service/internal/entity"
	"order-service/internal/repository"
	"order-service/internal/storage/cache"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	repo   repository.OrderRepository
	cache  *cache.Cache
}

func NewConsumer(brokers []string, topic, groupID string, repo repository.OrderRepository, cache *cache.Cache) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	return &Consumer{
		reader: r,
		repo:   repo,
		cache:  cache,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("[Kafka consumer] consumer is starting...")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("[Kafka consumer] error reading message: %v", err)
			continue
		}

		var order entity.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("[Kafka consumer] invalid message: %v", err)
			continue
		}

		if err := c.repo.Save(&order); err != nil {
			log.Printf("[Kafka consumer] failed to save order in DB: %v", err)
			continue
		}
		log.Printf("[Kafka consumer] saved order in DB: %s", order.OrderUID)

		c.cache.Set(&order)
		log.Printf("[Kafka consumer] saved order in cache: %s", order.OrderUID)
	}
}
