package repository

import (
	"AvitoWinter/internal/database"
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
