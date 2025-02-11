package repository

import (
	utils2 "AvitoWinter/internal/utils"
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

type Purchase struct {
	ID          int       `json:"id"`
	UserUUID    uuid.UUID `json:"user_uuid"`
	ItemUUID    uuid.UUID `json:"item_uuid"`
	Quantity    int       `json:"quantity"`
	TotalPrice  int       `json:"total_price"`
	DateCreated time.Time `json:"date_created"`
}

func NewPurchase(userUUID uuid.UUID, itemUUID uuid.UUID, quantity int, price int) *Purchase {
	serverTime := utils2.GetCurrentTime()
	purchase := &Purchase{UserUUID: userUUID, ItemUUID: itemUUID, Quantity: quantity, TotalPrice: quantity * price, DateCreated: serverTime}
	return purchase
}

type Ownership struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	ItemUUID uuid.UUID `json:"item_uuid"`
	Quantity int       `json:"quantity"`
}

func NewOwnership(userUUID uuid.UUID, itemUUID uuid.UUID, quantity int) *Ownership {
	return &Ownership{UserUUID: userUUID, ItemUUID: itemUUID, Quantity: quantity}
}

func (o *Ownership) IncQuantity() int {
	o.Quantity++
	return o.Quantity
}
