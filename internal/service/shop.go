package service

import (
	entity2 "AvitoWinter/internal/entity"
	"context"
)

type ShopRepo interface {
	GetInfoByUUID(ctx context.Context, userUUID string) error
	PutOrder(ctx context.Context)

	Ping() error
}

type ShopService struct {
	Repo ShopRepo
}

func NewShopService(repo ShopRepo) *ShopService {
	return &ShopService{Repo: repo}
}

func (s *ShopService) GetInfo(ctx context.Context) {

}

func (s *ShopService) TransferCoin(ctx context.Context, transferInfo entity2.TransferInfo) error {

	return nil
}

func (s *ShopService) PurchaseItem(ctx context.Context, purchasedInfo entity2.PurchaseInfo) error {

	return nil
}

func (s *ShopService) AuthenticationUser(ctx context.Context, userCredential entity2.UserCredentials) (string, error) {

	return "", nil
}
