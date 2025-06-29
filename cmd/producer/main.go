package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	e "order-service/internal/entity"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Kafka producer
	writer := kafka.Writer{
		Addr:     kafka.TCP("kafka:9092", "localhost:9092"),
		Topic:    "orders-topic",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	// 15 сообщений (заказов) для отправки
	for i := 0; i < 15; i++ {
		orderUID := fmt.Sprintf("orderTEST%d", rand.Intn(1000000))
		order := e.Order{
			OrderUID:    orderUID,
			TrackNumber: "KafkaTestPRODUCER",
			Entry:       "Kafka",
			Delivery: e.Delivery{
				Name:    "Test Testov",
				Phone:   "+9720000000",
				Zip:     "2639809",
				City:    "Kiryat Mozkin",
				Address: "Ploshad Mira 15",
				Region:  "Kraiot",
				Email:   "test@gmail.com",
			},
			Payment: e.Payment{
				Transaction:  "b563feb7b2b84b6test",
				RequestID:    "",
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       1817,
				PaymentDT:    1637907727,
				Bank:         "alpha",
				DeliveryCost: 1500,
				GoodsTotal:   317,
				CustomFee:    0,
			},
			Items: []e.Item{
				{
					ChrtID:      9934930,
					TrackNumber: "WBILMTESTTRACK",
					Price:       453,
					Rid:         "ab4219087a764ae0btest",
					Name:        "Mascaras",
					Sale:        30,
					Size:        "0",
					TotalPrice:  317,
					NmID:        2389212,
					Brand:       "Vivienne Sabo",
					Status:      202,
				},
			},
			Locale:            "en",
			InternalSignature: "",
			CustomerID:        "test",
			DeliveryService:   "meest",
			ShardKey:          "9",
			SmID:              99,
			DateCreated:       time.Now().Format(time.RFC3339),
			OofShard:          "1",
		}

		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("[Kafka producer] error marshalling order: %v\n", err)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: data,
		})
		if err != nil {
			log.Printf("[Kafka producer] failed to send message: %v\n", err)
		} else {
			log.Printf("[Kafka producer] Sent order: %s", order.OrderUID)
		}

		// Случайная задержка до 3 секунд
		delay := time.Duration(rand.Intn(3)) * time.Second
		time.Sleep(delay)
	}

	log.Println("[Kafka producer] Producer finished sending 10 messages.")
}
