package main

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	clusterID = "test-cluster"
	clientID  = "client-consumer"
	natsUrl   = "0.0.0.0:4223"
)

type Delivery struct {
	//Id      int    `db:"deliveries.id" json:"id"`
	Name    string `db:"name" json:"name,omitempty"`
	Phone   string `db:"phone" json:"phone,omitempty"`
	Zip     string `db:"zip" json:"zip,omitempty"`
	City    string `db:"city" json:"city,omitempty"`
	Address string `db:"address" json:"address,omitempty"`
	Region  string `db:"region" json:"region,omitempty"`
	Email   string `db:"email" json:"email,omitempty"`
}

type Payment struct {
	//Id           int    `db:"payments.id" json:"id"`
	Transaction  string `db:"transaction" json:"transaction,omitempty"`
	RequestId    string `db:"request_id" json:"request_id,omitempty"`
	Currency     string `db:"currency" json:"currency,omitempty"`
	Provider     string `db:"provider" json:"provider,omitempty"`
	Amount       int    `db:"amount" json:"amount,omitempty"`
	PaymentDt    int    `db:"payment_dt" json:"payment_dt,omitempty"`
	Bank         string `db:"bank" json:"bank,omitempty"`
	DeliveryCost int    `db:"delivery_cost" json:"delivery_cost,omitempty"`
	GoodsTotal   int    `db:"goods_total" json:"goods_total,omitempty"`
	CustomFee    int    `db:"custom_fee" json:"custom_fee,omitempty"`
}

type Item struct {
	Id          int    `db:"id" json:"id"`
	ChrtId      int    `db:"chrt_id" json:"chrt_id,omitempty"`
	TrackNumber string `db:"track_number" json:"track_number,omitempty"`
	Price       int    `db:"price" json:"price,omitempty"`
	Rid         string `db:"rid" json:"rid,omitempty"`
	Name        string `db:"name" json:"name,omitempty"`
	Sale        int    `db:"sale" json:"sale,omitempty"`
	Size        string `db:"size" json:"size,omitempty"`
	TotalPrice  int    `db:"total_price" json:"total_price,omitempty"`
	NmId        int    `db:"nm_id" json:"nm_id,omitempty"`
	Brand       string `db:"brand" json:"brand,omitempty"`
	Status      int    `db:"status" json:"status,omitempty"`
}

type Order struct {
	Id                int       `db:"id" json:"id"`
	OrderUID          string    `db:"order_uid" json:"order_uid,omitempty"`
	TrackNumber       string    `db:"track_number" json:"track_number,omitempty"`
	Entry             string    `db:"entry" json:"entry,omitempty"`
	Delivery          *Delivery `json:"delivery,omitempty"`
	Payment           *Payment  `json:"payment,omitempty"`
	Items             []Item    `json:"items"`
	Locale            string    `db:"locale" json:"locale,omitempty"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature"`
	CustomerId        string    `db:"customer_id" json:"customer_id,omitempty"`
	DeliveryService   string    `db:"delivery_service" json:"delivery_service,omitempty"`
	Shardkey          string    `db:"shardkey" json:"shardkey,omitempty"`
	SmId              int       `db:"sm_id" json:"sm_id,omitempty"`
	DateCreated       string    `db:"date_created" json:"date_created,omitempty"`
	OofShard          string    `db:"oof_shard" json:"oof_shard,omitempty"`
}

type OrderDB struct {
	Order
	Delivery
	Payment
	DeliveryId int `db:"delivery_id"`
	PaymentId  int `db:"payment_id"`
}

func (ordDB *OrderDB) ToOrder() Order {
	ord := ordDB.Order
	ord.Delivery = &ordDB.Delivery
	ord.Payment = &ordDB.Payment
	return ord
}

type ItemDB struct {
	Item
	OrderId int `db:"order_id"`
}

func (itmDB *ItemDB) ToItem() Item {
	return itmDB.Item
}

func seed(db *sqlx.DB) ([]Order, error) {
	var orders []Order
	var ordersDB []OrderDB

	err := db.Select(&ordersDB, "select orders.id, order_uid, track_number, entry, name, phone,"+
		" zip, city, address, region, email, transaction, request_id, currency, provider, amount, payment_dt,"+
		" bank, delivery_cost, goods_total, custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id,"+
		" date_created, oof_shard from orders join deliveries on orders.delivery_id = deliveries.id join payments on orders.payment_id = payments.id;")
	if err != nil {
		return nil, err
	}

	for _, ordDB := range ordersDB {
		var itemsDB []ItemDB
		err = db.Select(&itemsDB, "select * from items where order_id=$1;", ordDB.Id)
		if err != nil {
			return nil, err
		}
		var items []Item
		for _, itmDB := range itemsDB {
			items = append(items, itmDB.ToItem())
		}
		order := ordDB.ToOrder()
		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func storeDB(db *sqlx.DB, order *Order) error {
	tx, err := db.Begin()
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

func getOrder(res http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/orders/")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
	}

	var order Order
	for _, ord := range orders {
		if ord.Id == orderId {
			order = ord
			break
		}
	}

	tmpl, _ := template.ParseFiles("templates/order.html")
	tmpl.Execute(res, order)
}

var orders []Order

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}
	DBUser, _ := os.LookupEnv("DB_USER")
	DBPassword, _ := os.LookupEnv("DB_PASSWORD")
	DBName, _ := os.LookupEnv("DB_NAME")
	SSLMode, _ := os.LookupEnv("SSL_MODE")

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		DBUser, DBPassword, DBName, SSLMode))
	if err != nil {
		log.Println(err)
		return
	}
	orders, err = seed(db)
	if err != nil {
		log.Println(err)
		return
	}

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsUrl))
	if err != nil {
		log.Println(err)
		return
	}
	defer sc.Close()

	sub, err := sc.Subscribe("orders", func(m *stan.Msg) {
		var order Order
		log.Printf("Received a message: %s\n", string(m.Data))
		err = json.Unmarshal(m.Data, &order)
		if err != nil {
			log.Println(err)
			return
		}
		err = storeDB(db, &order)
		if err != nil {
			log.Println(err)
			return
		}
		orders = append(orders, order)
	})
	if err != nil {
		log.Println(err)
		return
	}

	defer sub.Unsubscribe()

	http.HandleFunc("/orders/", getOrder)
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
