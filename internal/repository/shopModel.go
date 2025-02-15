package repository

import (
	utils2 "AvitoWinter/internal/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Item struct {
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Coins    int    `json:"coins"`
}

func NewUser(username string, password string, coins int) *User {
	return &User{Username: username, Password: password, Coins: coins}
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
	Sender      string    `json:"sender"`
	Recipient   string    `json:"recipient"`
	Amount      int       `json:"amount"`
	DateCreated time.Time `json:"date_created"`
}

func NewTransfer(sender string, recipient string, amount int) *Transfer {
	return &Transfer{Sender: sender, Recipient: recipient, Amount: amount, DateCreated: utils2.GetCurrentTime()}
}

type Purchase struct {
	ID          int       `json:"id"`
	User        string    `json:"user_uuid"`
	Item        string    `json:"item_uuid"`
	Quantity    int       `json:"quantity"`
	TotalPrice  int       `json:"total_price"`
	DateCreated time.Time `json:"date_created"`
}

func NewPurchase(user string, item string, quantity int, price int) *Purchase {
	serverTime := utils2.GetCurrentTime()
	purchase := &Purchase{User: user, Item: item, Quantity: quantity, TotalPrice: quantity * price, DateCreated: serverTime}
	return purchase
}

type Ownership struct {
	User     string `json:"user"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

func NewOwnership(user string, item string, quantity int) *Ownership {
	return &Ownership{User: user, Item: item, Quantity: quantity}
}

func (o *Ownership) IncQuantity() int {
	o.Quantity++
	return o.Quantity
}
