package main

import (
	"context"
	"log"
	"net/http"
	"order-service/config"
	"order-service/internal/handlers"
	"order-service/internal/repository"
	"order-service/internal/storage"
	"order-service/internal/storage/cache"
	"time"

	"order-service/internal/kafka"

	"github.com/gorilla/mux"
)

const (
	cacheSize int = 5 // Размер кэша (количество хранимых заказов)
	cacheTTL  int = 5 // Количество минут - TTL для кэша
)

func main() {
	cfg := config.LoadConfig()
	db := storage.ConnectDB(cfg)
	defer db.Close()
	c := cache.NewCache(cacheSize, time.Duration(cacheTTL)*time.Minute)

	orderRepo := repository.NewPostgresOrderRepository(db)
	orderHandler := handlers.NewHandler(c, orderRepo)

	// Восстановление кэша
	orders, err := orderRepo.GetAll(cacheSize)
	if err != nil {
		log.Fatalf("[error] failed to preload cache: %v", err)
	}
	c.Load(orders)
	log.Printf("Cache preloaded...\n")

	router := mux.NewRouter()
	router.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET")
	router.Handle("/", http.FileServer(http.Dir("web")))

	// HTTP сервер
	go func() {
		log.Println("Server is starting on: 8081...")
		log.Fatal(http.ListenAndServe(":8081", router))
	}()

	// Kafka consumer
	consumer := kafka.NewConsumer(
		[]string{"kafka:9092"},
		"orders-topic",
		"orders-group",
		orderRepo,
		c,
	)
	ctx := context.Background()
	consumer.Start(ctx)
}
