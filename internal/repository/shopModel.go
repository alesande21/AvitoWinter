package repository

import (
	utils2 "AvitoWinter/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func NewUser(UUID uuid.UUID, username string, password string, coins int) *User {
	return &User{UUID: UUID, Username: username, Password: password, Coins: coins}
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return fmt.Errorf("-> bcrypt.GenerateFromPassword%v", err)
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
	if err != nil {
		return fmt.Errorf("-> bcrypt.CompareHashAndPassword%v", err)
	}
	return nil
}

type Transfer struct {
	ID          int       `json:"id"`
	Sender      uuid.UUID `json:"sender"`
	Recipient   uuid.UUID `json:"recipient"`
	Amount      int       `json:"amount"`
	DateCreated time.Time `json:"date_created"`
}

func NewTransfer(sender uuid.UUID, recipient uuid.UUID, amount int) *Transfer {
	return &Transfer{Sender: sender, Recipient: recipient, Amount: amount, DateCreated: utils2.GetCurrentTime()}
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
