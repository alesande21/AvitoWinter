package repository

import (
	"AvitoWinter/internal/database"
	entity2 "AvitoWinter/internal/entity"
	"context"
	"fmt"
	"log"
)

//type Cache interface {
//	Set(ctx context.Context, key string, user *entity2.UserInfo, expiration time.Duration) error
//	Get(ctx context.Context, key string) (*entity2.UserInfo, error)
//}
//
//type ShopRepoWithCache struct {
//	shopRepo *ShopRepo
//	cache    Cache
//}
//
//func NewShopRepoWithCache(userRepo *ShopRepo, cache Cache) *ShopRepoWithCache {
//	return &ShopRepoWithCache{shopRepo: userRepo, cache: cache}
//}

type ShopRepo struct {
	dbRepo database.DBRepository
}

func NewShopRepo(dbRepo database.DBRepository) *ShopRepo {
	return &ShopRepo{dbRepo: dbRepo}
}

func (s ShopRepo) GetInfoByUUID(ctx context.Context, username string) (*entity2.UserInfo, error) {
	user, err := s.getUserByUseUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("-> s.getUserByUseUUID%v", err)
	}

	items, err := s.getUserItemsByUserUUID(ctx, user.UUID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, i.product_name, quantity
		FROM transfers t
		LEFT JOIN users u ON t.sender = u.uuid
		WHERE user_uuid = $1
	`

	fmt.Println(query, items)

	return nil, nil
}

func (s ShopRepo) CheckUser(ctx context.Context, userCredential entity2.UserCredentials) (string, error) {
	repoCredential, err := s.getUserByUsername(ctx, userCredential.Password())
	if err != nil {
		return "", fmt.Errorf("-> r.dbRepo.QueryRow.Scan: пользователь по идентификатору %s не найден: %w", userCredential.Identifier(), err)
	}

	err = repoCredential.CheckPassword(userCredential.Password())
	if err != nil {
		return "", fmt.Errorf("-> repoCredential.CheckPassword%v", err)
	}

	return repoCredential.Username, nil
}

func (s ShopRepo) PutPurchaseInfo(ctx context.Context, info entity2.PurchaseInfo) error {
	queryInsertPurchase := `
		INSERT INTO purchases (username, item, quantity, total_price, date_created)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, item, quantity, total_price, date_created
	`

	queryInsertOwnership := `
		INSERT INTO ownership (username, item, quantity)
		VALUES ($1, $2, $3)
		RETURNING username, item, quantity
	`

	queryUpdateOwnership := `
		UPDATE ownership
		SET quantity = $3
		WHERE username = $1 AND item = $2
		RETURNING username, item, quantity
	`

	queryUpdateCoins := `
		UPDATE users
		SET coins = $2
		WHERE username = $1
		RETURNING coins
	`

	item, err := s.getItemByProductName(ctx, info.Item())
	if err != nil {
		return fmt.Errorf("-> s.getItemByProductName%v", err)
	}

	//TODO так как есть аутификация возможно это не нужно, хотя нужно для нахождения UUID юзера
	user, err := s.getUserByUsername(ctx, info.Username())
	if err != nil {
		return fmt.Errorf("-> s.getUserByUsername%v", err)
	}

	if user.Coins < item.Price {
		return fmt.Errorf(": недостаточно монет на счете. Монет - %d, необходимо - %d", user.Coins, item.Price)
	}

	repoPurchase := NewPurchase(user.Username, item.ProductName, info.Quantity(), item.Price)

	tx, err := s.dbRepo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("-> r.dbRepo.BeginTx: не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, queryInsertPurchase, repoPurchase.User, repoPurchase.Item,
		repoPurchase.Quantity, repoPurchase.TotalPrice, repoPurchase.DateCreated)
	err = row.Scan(&repoPurchase.ID, &repoPurchase.User, &repoPurchase.Item, &repoPurchase.Quantity,
		&repoPurchase.TotalPrice, &repoPurchase.DateCreated)
	if err != nil {
		log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
		return fmt.Errorf("-> row.Scan:%s", err)
	}

	userOwnership, err := s.getOwnershipByUserAndItem(ctx, user.Username, item.ProductName)
	if err != nil {
		row = tx.QueryRowContext(ctx, queryInsertOwnership, user.Username, item.ProductName, 1)
		err = row.Scan(&userOwnership.User, &userOwnership.Item, &userOwnership.Quantity)
		if err != nil {
			log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
			return fmt.Errorf("-> row.Scan:%s", err)
		}
	} else {
		row = tx.QueryRowContext(ctx, queryUpdateOwnership, userOwnership.User, userOwnership.Item, userOwnership.IncQuantity())
		err = row.Scan(&userOwnership.User, &userOwnership.Item, &userOwnership.Quantity)
		if err != nil {
			log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
			return fmt.Errorf("-> row.Scan:%s", err)
		}
	}

	row = tx.QueryRowContext(ctx, queryUpdateCoins, user.Username, user.Coins-item.Price)
	err = row.Scan(&user.Coins)
	if err != nil {
		log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
		return fmt.Errorf("-> row.Scan:%s", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("-> tx.Commit: не удалось завершить транзакцию: %w", err)
	}

	return nil
}

func (s ShopRepo) PutTransferInfo(ctx context.Context, info entity2.TransferInfo) error {
	sender, err := s.getUserByUsername(ctx, info.Sender())
	if err != nil {
		return fmt.Errorf("-> s.getUserByUseUUID%v", err)
	}

	recipient, err := s.getUserByUsername(ctx, info.Recipient())
	if err != nil {
		return fmt.Errorf("-> s.getUserByUsername%v", err)
	}

	if sender.Coins < info.Amount() {
		return fmt.Errorf(": недостаточно средств для перевода")
	}

	queryInsert := `
		INSERT INTO transfers (sender, recipient, amount, date_created)
		VALUES ($1, $2, $3, $4)
		RETURNING id, sender, recipient, amount, date_created
	`

	queryUpdate := `
		WITH 
		sender_update AS (
			UPDATE users 
			SET coins = $1
			WHERE username = $2 
			RETURNING username, coins
		),
		recipient_update AS (
			UPDATE users 
			SET coins = $3
			WHERE username = $4 
			RETURNING username, coins
		)
		SELECT * FROM sender_update, recipient_update;
	`

	repoTransfer := NewTransfer(sender.Username, recipient.Username, info.Amount())

	tx, err := s.dbRepo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("-> r.dbRepo.BeginTx: не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, queryInsert, repoTransfer.Sender, repoTransfer.Recipient, repoTransfer.Amount, repoTransfer.DateCreated)
	err = row.Scan(&repoTransfer.ID, &repoTransfer.Sender, &repoTransfer.Recipient, &repoTransfer.Amount, &repoTransfer.DateCreated)
	if err != nil {
		log.Printf("Ошибка выполнения запроса в PutTransferInfo: %v\n", err)
		return fmt.Errorf("-> row.Scan:%s", err)
	}

	row = tx.QueryRowContext(ctx, queryUpdate, sender.Coins-repoTransfer.Amount, sender.Username, recipient.Username,
		recipient.Coins+repoTransfer.Amount)
	err = row.Scan(&sender.Username, &sender.Coins, &recipient.Username, &recipient.Coins)
	if err != nil {
		log.Printf("Ошибка выполнения запроса в PutTransferInfo: %v\n", err)
		return fmt.Errorf("-> row.Scan:%s", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("-> tx.Commit: не удалось завершить транзакцию: %w", err)
	}

	return nil
}

func (s ShopRepo) Ping() error {
	return nil
}

func (s ShopRepo) getItemByProductName(ctx context.Context, productName string) (*Item, error) {
	query := `
		SELECT product_name, price
		FROM items
		WHERE product_name = $1
	`

	var item *Item

	row := s.dbRepo.QueryRow(ctx, query, productName)
	err := row.Scan(&item.ProductName, &item.Price)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: продукт по productName %s не найден: %w", productName, err)
	}

	return item, nil
}

func (s ShopRepo) getUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT username, password, coins
		FROM users
		WHERE username = $1
	`

	var user *User

	row := s.dbRepo.QueryRow(ctx, query, username)
	err := row.Scan(&user.Username, &user.Password, &user.Coins)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: пользователь по username %s не найден: %w", username, err)
	}

	return user, nil
}

func (s ShopRepo) getOwnershipByUserAndItem(ctx context.Context, user string, item string) (*Ownership, error) {
	query := `
		SELECT username, item, quantity
		FROM ownership
		WHERE username = $1 AND item = $2
	`

	var own *Ownership

	row := s.dbRepo.QueryRow(ctx, query, user, item)
	err := row.Scan(&own.User, &own.Item, &own.Quantity)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: владение по user - %s и product - %s: %w", own.User, own.Item, err)
	}

	return own, nil
}

func (s ShopRepo) getUserItemsByUsername(ctx context.Context, username string) ([]Ownership, error) {
	query := `
		SELECT user_uuid, i.product_name, quantity
		FROM ownership o
		LEFT JOIN items i ON o.items_uuid = i.uuid
		WHERE user_uuid = $1
	`

	rows, err := s.dbRepo.Query(ctx, query, userUUID)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var userItems []Ownership
	for rows.Next() {
		var userItem Ownership
		err = rows.Scan(&userItem.UserUUID, &userItem.ItemUUID)
		if err != nil {
			log.Printf("ошибка выполнения: %v\n", err)
			return nil, err
		}
		userItems = append(userItems, userItem)
	}

	return userItems, nil
}
