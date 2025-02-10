package repository

import (
	"github.com/google/uuid"
	"time"
)

type Item struct {
	UUID        uuid.UUID `json:"uuid"`
	ProductName string    `json:"product_name"`
	Price       int       `json:"price"`
}

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	Coins    int       `json:"coins"`
}

type Transfer struct {
	ID          int       `json:"id"`
	Sender      uuid.UUID `json:"sender"`
	Recipient   uuid.UUID `json:"recipient"`
	Amount      int       `json:"amount"`
	DateCreated time.Time `json:"date_created"`
}

type Order struct {
	ID          int       `json:"id"`
	UserUUID    uuid.UUID `json:"user_uuid"`
	ItemUUID    uuid.UUID `json:"item_uuid"`
	Quantity    int       `json:"quantity"`
	TotalPrice  int       `json:"total_price"`
	DateCreated time.Time `json:"date_created"`
}

type Ownership struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	ItemUUID uuid.UUID `json:"item_uuid"`
	Quantity int       `json:"quantity"`
}
