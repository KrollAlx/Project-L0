package models

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
