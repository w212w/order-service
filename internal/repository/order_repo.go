package repository

import (
	"database/sql"
	"fmt"
	"order-service/internal/entity"
	"time"
)

type OrderRepository interface {
	Save(order *entity.Order) error
	GetByUID(orderUID string) (*entity.Order, error)
	GetAll(cacheSize int) ([]*entity.Order, error)
}

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Save(order *entity.Order) error {
	// Начало транзакции
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Вставка заказа
	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service,
			shardkey, sm_id, date_created, oof_shard
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) `, order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard)

	if err != nil {
		return fmt.Errorf("[error] insert order: %w", err)
	}

	// Вставка доставки
	_, err = tx.Exec(`
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)

	if err != nil {
		return fmt.Errorf("[error] insert delivery: %w", err)
	}

	// Вставка оплаты
	_, err = tx.Exec(`
		INSERT INTO payment (
			order_uid, transaction, request_id, currency,
			provider, amount, payment_dt, bank, delivery_cost,
			goods_total, custom_fee
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDT, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee)

	if err != nil {
		return fmt.Errorf("[error] insert payment: %w", err)
	}

	// Вставка items
	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid,
				name, sale, size, total_price, nm_id, brand, status
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status)

		if err != nil {
			return fmt.Errorf("[error] insert item: %w", err)
		}
	}
	// Коммит транзакции
	return tx.Commit()
}

func (r *PostgresOrderRepository) GetByUID(uid string) (*entity.Order, error) {
	order := &entity.Order{}
	delivery := entity.Delivery{}
	payment := entity.Payment{}
	items := []entity.Item{}

	// ЗАДЕРЖКА для проверки работы кэша
	time.Sleep(1 * time.Second)
	// Поиск заказа в БД
	err := r.db.QueryRow(`
		SELECT order_uid, track_number, entry, locale, internal_signature,
		       customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
	`, uid).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry,
		&order.Locale, &order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.ShardKey, &order.SmID,
		&order.DateCreated, &order.OofShard)

	if err != nil {
		return nil, fmt.Errorf("[error] order not found: %w", err)
	}

	// Поиск доставки в БД
	err = r.db.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uid = $1
	`, uid).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip,
		&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("[error] delivery not found: %w", err)
	}

	// Поиск оплаты в БД
	err = r.db.QueryRow(`
		SELECT transaction, request_id, currency, provider, amount,
		       payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid = $1
	`, uid).Scan(&payment.Transaction, &payment.RequestID, &payment.Currency,
		&payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank,
		&payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("[error] payment not found: %w", err)
	}

	// Поиск товаров в БД
	rows, err := r.db.Query(`
		SELECT chrt_id, track_number, price, rid, name, sale, size,
		       total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
	`, uid)
	if err != nil {
		return nil, fmt.Errorf("[error] items not found: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("[error] item scan: %w", err)
		}
		items = append(items, item)
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	return order, nil
}

func (r *PostgresOrderRepository) GetAll(cacheSize int) ([]*entity.Order, error) {
	// Поиск всех uid заказов в БД
	rows, err := r.db.Query(fmt.Sprintf(`SELECT order_uid FROM orders ORDER BY date_created DESC LIMIT %d`, cacheSize))
	if err != nil {
		return nil, fmt.Errorf("[error] get all orders: %w", err)
	}
	defer rows.Close()

	var orders []*entity.Order

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			continue
		}

		order, err := r.GetByUID(uid)
		if err != nil {
			continue
		}

		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[error] row iteration failed: %w", err)
	}

	return orders, nil
}
