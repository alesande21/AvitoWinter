package entity

import (
	"fmt"
	"regexp"
)

type Item struct {
	item     string
	quantity int
}

func NewItem(item string, quantity int) *Item {
	return &Item{
		item:     item,
		quantity: quantity,
	}
}

func (i *Item) GetItem() string {
	return i.item
}

func (i *Item) GetQuantity() int {
	return i.quantity
}

type Transfer struct {
	username string
	amount   int
}

func NewTransfer(username string, amount int) *Transfer {
	return &Transfer{
		username: username,
		amount:   amount,
	}
}

func (t *Transfer) GetUsername() string {
	return t.username
}

func (t *Transfer) GetAmount() int {
	return t.amount
}

type UserTransfers struct {
	received []*Transfer
	sent     []*Transfer
}

func NewUserTransfers(received, sent []*Transfer) *UserTransfers {
	return &UserTransfers{
		received: received,
		sent:     sent,
	}
}

func (ut *UserTransfers) GetReceived() []*Transfer {
	return ut.received
}

func (ut *UserTransfers) GetSent() []*Transfer {
	return ut.sent
}

type UserInfo struct {
	coins         int
	items         []*Item
	userTransfers *UserTransfers
}

func NewUserInfo(coins int, items []*Item, userTransfers *UserTransfers) *UserInfo {
	return &UserInfo{
		coins:         coins,
		items:         items,
		userTransfers: userTransfers,
	}
}

func (ui *UserInfo) GetCoins() *int {
	return &ui.coins
}

func (ui *UserInfo) GetItems() []*Item {
	return ui.items
}

func (ui *UserInfo) GetUserTransfers() *UserTransfers {
	return ui.userTransfers
}

type PurchaseInfo struct {
	username string
	item     string
	quantity int
}

func (p PurchaseInfo) Username() string {
	return p.username
}

func (p PurchaseInfo) Item() string {
	return p.item
}

func (p PurchaseInfo) Quantity() int {
	return p.quantity
}

func NewPurchaseInfo(username, item string, quantity int) (*PurchaseInfo, error) {
	if item == "" || username == "" {
		return nil, fmt.Errorf(": не задан предмет покупки: %s", item)
	}
	return &PurchaseInfo{username: username, item: item, quantity: quantity}, nil
}

type TransferInfo struct {
	sender    string
	recipient string
	amount    int
}

func (t TransferInfo) Sender() string {
	return t.sender
}

func (t TransferInfo) Recipient() string {
	return t.recipient
}

func (t TransferInfo) Amount() int {
	return t.amount
}

func NewTransferInfo(sender string, recipient string, amount int) (*TransferInfo, error) {
	ti := &TransferInfo{sender: sender, recipient: recipient, amount: amount}
	err := ti.validate()
	if err != nil {
		return nil, fmt.Errorf("-> ti.validate%v", err)
	}

	return ti, nil
}

func (t TransferInfo) validate() error {
	err := validateIdentifier(t.sender)
	if err != nil {
		return fmt.Errorf("-> validateIdentifier%v", err)
	}

	err = validateIdentifier(t.recipient)
	if err != nil {
		return fmt.Errorf("-> validateIdentifier%v", err)
	}

	if t.amount <= 0 {
		return fmt.Errorf(": сумма не может быть меньше отрицательной или меньше нуля: %d", t.amount)
	}

	return nil
}

type UserCredentials struct {
	identifier string
	password   string
}

func (ru UserCredentials) Identifier() string {
	return ru.identifier
}

func (ru UserCredentials) Password() string {
	return ru.password
}

func NewUserCredentials(identifier string, password string) (*UserCredentials, error) {
	ru := &UserCredentials{identifier: identifier, password: password}
	err := ru.validate()

	if err != nil {
		return nil, fmt.Errorf("-> ru.validate%v", err)
	}

	return &UserCredentials{identifier: identifier, password: password}, nil
}

func (ru UserCredentials) validate() error {
	err := validateIdentifier(ru.identifier)
	if err != nil {
		return fmt.Errorf("-> validateIdentifier%v", err)
	}

	err = validatePassword(ru.password)
	if err != nil {
		return fmt.Errorf("-> validatePassword%v", err)
	}

	return nil
}

func validateIdentifier(identifier string) error {
	identifierRegex := `^[a-zA-Z0-9._-]{4,50}$`

	matched, err := regexp.MatchString(identifierRegex, identifier)
	if err != nil {
		return fmt.Errorf("ошибка при проверке идентификатора: %v", err)
	}

	if !matched {
		return fmt.Errorf("идентификатор должен быть длиной от 4 до 50 символов и содержать только буквы, цифры, точки, дефисы и символы подчеркивания")
	}

	return nil
}

func validatePassword(password string) error {
	passwordRegex := `^.{4,}$`

	matched, err := regexp.MatchString(passwordRegex, password)
	if err != nil {
		return fmt.Errorf("ошибка при проверке пароля: %v", err)
	}

	if !matched {
		return fmt.Errorf("пароль должен быть длиной не менее 4 символов")
	}

	return nil
}
