package service

import (
	auth2 "AvitoWinter/internal/auth"
	entity2 "AvitoWinter/internal/entity"
	"context"
	"fmt"
)

type ShopRepo interface {
	GetInfoByUUID(ctx context.Context, userUUID string) error
	PutPurchaseInfo(ctx context.Context, info entity2.PurchaseInfo) error
	CheckUser(ctx context.Context, userCredential entity2.UserCredentials) (string, error)
	PutTransferInfo(ctx context.Context, info entity2.TransferInfo) error

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
	s.repo.PutTransferInfo(ctx, transferInfo)
	return nil
}

func (s *ShopService) PurchaseItem(ctx context.Context, info entity2.PurchaseInfo) error {
	return s.PurchaseItem(ctx, info)
}

func (s *ShopService) AuthenticationUser(ctx context.Context, userCredential entity2.UserCredentials) (string, error) {
	UUID, err := s.repo.CheckUser(ctx, userCredential)
	if err != nil {
		return "", fmt.Errorf("-> s.Repo.CheckUser%v", err)
	}

	tokenString, err := auth2.GenerateJWT(UUID)
	if err != nil {
		return "", fmt.Errorf("-> auth2.GenerateJWT%v", err)
	}

	return tokenString, nil
}

func (s *ShopService) ValidateToken(ctx context.Context, tokenString string) (*auth2.JWTClaim, error) {
	return auth2.ValidateToken(tokenString)
}
