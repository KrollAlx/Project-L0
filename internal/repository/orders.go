package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"project-L0/internal/models"
)

type Orders interface {
	GetAll() ([]models.Order, error)
	Create(order *models.Order) error
}

type OrderDB struct {
	models.Order
	models.Delivery
	models.Payment
	DeliveryId int `db:"delivery_id"`
	PaymentId  int `db:"payment_id"`
}

func (ordDB *OrderDB) ToOrder() models.Order {
	ord := ordDB.Order
	ord.Delivery = &ordDB.Delivery
	ord.Payment = &ordDB.Payment
	return ord
}

type ItemDB struct {
	models.Item
	OrderId int `db:"order_id"`
}

func (itmDB *ItemDB) ToItem() models.Item {
	return itmDB.Item
}

type OrdersRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{db: db}
}

func (r *OrdersRepo) GetAll() ([]models.Order, error) {
	var orders []models.Order
	var ordersDB []OrderDB

	err := r.db.Select(&ordersDB, "select orders.id, order_uid, track_number, entry, name, phone,"+
		" zip, city, address, region, email, transaction, request_id, currency, provider, amount, payment_dt,"+
		" bank, delivery_cost, goods_total, custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id,"+
		" date_created, oof_shard from orders join deliveries on orders.delivery_id = deliveries.id join payments on orders.payment_id = payments.id;")
	if err != nil {
		return nil, err
	}

	for _, ordDB := range ordersDB {
		var itemsDB []ItemDB
		err = r.db.Select(&itemsDB, "select * from items where order_id=$1;", ordDB.Id)
		if err != nil {
			return nil, err
		}
		var items []models.Item
		for _, itmDB := range itemsDB {
			items = append(items, itmDB.ToItem())
		}
		order := ordDB.ToOrder()
		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrdersRepo) Create(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	insertPaymentQuery := "insert into payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) values " +
		"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;"
	row := tx.QueryRow(insertPaymentQuery, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	var paymentId int
	err = row.Scan(&paymentId)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertDeliveryQuery := "insert into deliveries (name, phone, zip, city, address, region, email) " +
		"values ($1, $2, $3, $4, $5, $6, $7) returning id;"
	row = tx.QueryRow(insertDeliveryQuery, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	var deliveryId int
	err = row.Scan(&deliveryId)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertOrderQuery := "insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard, delivery_id, payment_id)" +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id;"
	row = tx.QueryRow(insertOrderQuery, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard, deliveryId, paymentId)
	var orderId int
	err = row.Scan(&orderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertItemQuery := "insert into items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id) values"
	for _, item := range order.Items {
		_, err = tx.Exec(insertItemQuery+"($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);",
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status, orderId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
