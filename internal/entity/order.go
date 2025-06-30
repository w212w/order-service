package entity

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shard_key"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type OrderResponse struct {
	OrderUID    string           `json:"order_uid"`
	TrackNumber string           `json:"track_number"`
	Entry       string           `json:"entry"`
	Delivery    DeliveryResponse `json:"delivery"`
	Payment     PaymentResponse  `json:"payment"`
	Items       []ItemResponse   `json:"items"`
	DateCreated string           `json:"date_created"`
}

type DeliveryResponse struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type PaymentResponse struct {
	Amount       int    `json:"amount"`
	Currency     string `json:"currency"`
	GoodsTotal   int    `json:"goods_total"`
	DeliveryCost int    `json:"delivery_cost"`
}

type ItemResponse struct {
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Brand      string `json:"brand"`
	TotalPrice int    `json:"total_price"`
}
