package service

import (
	entity2 "AvitoWinter/internal/entity"
	"context"
)

type ShopRepo interface {
	GetInfoByUUID(ctx context.Context, userUUID string) error
	PutPurchaseInfo(ctx context.Context, info entity2.PurchaseInfo) error

	Ping() error
}

type ShopService struct {
	repo ShopRepo
}

func NewShopService(repo ShopRepo) *ShopService {
	return &ShopService{repo: repo}
}

func (s *ShopService) GetInfo(ctx context.Context) {

}

func (s *ShopService) TransferCoin(ctx context.Context, transferInfo entity2.TransferInfo) error {

	return nil
}

func (s *ShopService) PurchaseItem(ctx context.Context, info entity2.PurchaseInfo) error {
	return s.PurchaseItem(ctx, info)
}

func (s *ShopService) AuthenticationUser(ctx context.Context, userCredential entity2.UserCredentials) (string, error) {

	return "", nil
}
