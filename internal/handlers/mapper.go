package handlers

import e "order-service/internal/entity"

func ConvertToResponse(o *e.Order) *e.OrderResponse {
	items := make([]e.ItemResponse, 0, len(o.Items))
	for _, it := range o.Items {
		items = append(items, e.ItemResponse{
			Name:       it.Name,
			Price:      it.Price,
			Brand:      it.Brand,
			TotalPrice: it.TotalPrice,
		})
	}

	return &e.OrderResponse{
		OrderUID:    o.OrderUID,
		TrackNumber: o.TrackNumber,
		Entry:       o.Entry,
		Delivery: e.DeliveryResponse{
			Name:    o.Delivery.Name,
			Phone:   o.Delivery.Phone,
			City:    o.Delivery.City,
			Address: o.Delivery.Address,
		},
		Payment: e.PaymentResponse{
			Amount:       o.Payment.Amount,
			Currency:     o.Payment.Currency,
			GoodsTotal:   o.Payment.GoodsTotal,
			DeliveryCost: o.Payment.DeliveryCost,
		},
		Items:       items,
		DateCreated: o.DateCreated,
	}
}
