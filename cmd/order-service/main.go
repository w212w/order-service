package main

import (
	"order-service/config"
	"order-service/internal/storage"
	"order-service/internal/storage/cache"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	db := storage.ConnectDB(cfg)
	defer db.Close()
	c := cache.NewCache(100, 1*time.Minute)

	// data, err := os.ReadFile("api/model.json")
	// if err != nil {
	// 	log.Fatalf("failed to read JSON: %v", err)
	// }

	// var order entity.Order
	// if err := json.Unmarshal(data, &order); err != nil {
	// 	log.Fatalf("failed to parse JSON: %v", err)
	// }

	// // Сохраняем заказ
	// repo := repository.NewPostgresOrderRepository(db)
	// if err := repo.Save(&order); err != nil {
	// 	log.Fatalf("failed to save order: %v", err)
	// }

	// fmt.Println("Order saved successfully!")
}
