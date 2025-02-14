package repository

import (
	"AvitoWinter/internal/database"
	entity2 "AvitoWinter/internal/entity"
	"context"
	"fmt"
	"github.com/google/uuid"
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

func (s ShopRepo) GetInfoByUUID(ctx context.Context, userUUID string) error {
	//TODO implement me
	panic("implement me")
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

	return repoCredential.UUID.String(), nil
}

func (s ShopRepo) PutPurchaseInfo(ctx context.Context, info entity2.PurchaseInfo) error {
	queryInsertPurchase := `
		INSERT INTO purchases (user_uuid, items_uuid, quantity, total_price, date_created)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_uuid, items_uuid, quantity, total_price, date_created
	`

	queryInsertOwnership := `
		INSERT INTO purchases (user_uuid, items_uuid, quantity)
		VALUES ($1, $2, $3)
		RETURNING user_uuid, items_uuid, quantity
	`

	queryUpdateOwnership := `
		UPDATE ownership
		SET quantity = $3
		WHERE user_uuid = $1 AND items_uuid = $2
		RETURNING user_uuid, items_uuid, quantity
	`

	queryUpdateCoins := `
		UPDATE users
		SET coins = $2
		WHERE user_uuid = $1
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

	repoPurchase := NewPurchase(user.UUID, item.UUID, info.Quantity(), item.Price)

	tx, err := s.dbRepo.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("-> r.dbRepo.BeginTx: не удалось начать транзакцию: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.QueryRowContext(ctx, queryInsertPurchase, repoPurchase.UserUUID, repoPurchase.ItemUUID,
		repoPurchase.Quantity, repoPurchase.TotalPrice, repoPurchase.DateCreated)
	err = row.Scan(&repoPurchase.ID, &repoPurchase.UserUUID, &repoPurchase.ItemUUID, &repoPurchase.Quantity,
		&repoPurchase.TotalPrice, &repoPurchase.DateCreated)
	if err != nil {
		log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
		return fmt.Errorf("-> row.Scan:%s", err)
	}

	userOwnership, err := s.getOwnershipByProduct(ctx, user.UUID, item.UUID)
	if err != nil {
		row = tx.QueryRowContext(ctx, queryInsertOwnership, user.UUID, item.UUID, 1)
		err = row.Scan(&userOwnership.UserUUID, &userOwnership.ItemUUID, &userOwnership.Quantity)
		if err != nil {
			log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
			return fmt.Errorf("-> row.Scan:%s", err)
		}
	} else {
		row = tx.QueryRowContext(ctx, queryUpdateOwnership, userOwnership.UserUUID, userOwnership.ItemUUID, userOwnership.IncQuantity())
		err = row.Scan(&userOwnership.UserUUID, &userOwnership.ItemUUID, &userOwnership.Quantity)
		if err != nil {
			log.Printf("Ошибка выполнения запроса в PutPurchaseInfo: %v\n", err)
			return fmt.Errorf("-> row.Scan:%s", err)
		}
	}

	row = tx.QueryRowContext(ctx, queryUpdateCoins, user.UUID, user.Coins-item.Price)
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
	recipient, err := s.getUserByUsername(ctx, info.RecipientUsername())

	return nil
}

func (s ShopRepo) Ping() error {
	//TODO implement me
	panic("implement me")
}

func (s ShopRepo) getItemByProductName(ctx context.Context, productName string) (*Item, error) {
	query := `
		SELECT uuid, product_name, price
		FROM items
		WHERE product_name = $1
	`

	var item *Item

	row := s.dbRepo.QueryRow(ctx, query, productName)
	err := row.Scan(&item.UUID, &item.ProductName, &item.Price)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: продукт по productName %s не найден: %w", productName, err)
	}

	return item, nil
}

func (s ShopRepo) getUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT uuid, username, password, coins
		FROM users
		WHERE username = $1
	`

	var user *User

	row := s.dbRepo.QueryRow(ctx, query, username)
	err := row.Scan(&user.UUID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: пользователь по username %s не найден: %w", username, err)
	}

	return user, nil
}

func (s ShopRepo) getUserByUseUUID(ctx context.Context, userUUID string) (*User, error) {
	query := `
		SELECT uuid, username, password, coins
		FROM users
		WHERE uuid = $1
	`

	var user *User

	row := s.dbRepo.QueryRow(ctx, query, userUUID)
	err := row.Scan(&user.UUID, &user.Username, &user.Password, &user.Coins)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: пользователь по userUUID %s не найден: %w", userUUID, err)
	}

	return user, nil
}

func (s ShopRepo) getOwnershipByProduct(ctx context.Context, userUUID uuid.UUID, productUUID uuid.UUID) (*Ownership, error) {
	query := `
		SELECT user_uuid, items_uuid, quantity
		FROM ownership
		WHERE user_uuid = $1 AND items_uuid
	`

	var own *Ownership

	row := s.dbRepo.QueryRow(ctx, query, userUUID, productUUID)
	err := row.Scan(&own.UserUUID, &own.ItemUUID, &own.Quantity)
	if err != nil {
		return nil, fmt.Errorf("-> r.dbRepo.QueryRow.Scan: владение по userUUID - %s и productUUID - %s: %w", own.UserUUID, own.ItemUUID, err)
	}

	return own, nil
}
