package main

import (
	"log"
	"net/http"
	"order-service/config"
	"order-service/internal/handlers"
	"order-service/internal/repository"
	"order-service/internal/storage"
	"order-service/internal/storage/cache"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	db := storage.ConnectDB(cfg)
	defer db.Close()
	c := cache.NewCache(100, 1*time.Minute)

	orderRepo := repository.NewPostgresOrderRepository(db)
	orderHandler := handlers.NewHandler(c, orderRepo)

	router := mux.NewRouter()
	router.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET")

	log.Println("Server is starting on: 8081...")
	log.Fatal(http.ListenAndServe(":8081", router))

}
