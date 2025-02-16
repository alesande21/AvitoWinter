package service

import (
	entity2 "AvitoWinter/internal/entity"
	"context"
	"fmt"
)

type ShopRepo interface {
	GetInfo(ctx context.Context, username string) (*entity2.UserInfo, error)
	PutPurchaseInfo(ctx context.Context, info entity2.PurchaseInfo) error
	CheckUser(ctx context.Context, userCredential *entity2.UserCredentials) (string, error)
	PutUser(ctx context.Context, userCredential *entity2.UserCredentials) (string, error)
	PutTransferInfo(ctx context.Context, info entity2.TransferInfo) error
}

type ShopService struct {
	repo ShopRepo
}

func NewShopService(repo ShopRepo) *ShopService {
	return &ShopService{repo: repo}
}

func (s *ShopService) GetInfo(ctx context.Context, username string) (*entity2.UserInfo, error) {
	return s.repo.GetInfo(ctx, username)
}

func (s *ShopService) TransferCoin(ctx context.Context, transferInfo entity2.TransferInfo) error {
	return s.repo.PutTransferInfo(ctx, transferInfo)
}

func (s *ShopService) PurchaseItem(ctx context.Context, info entity2.PurchaseInfo) error {
	return s.repo.PutPurchaseInfo(ctx, info)
}

func (s *ShopService) GetUserByCredentials(ctx context.Context, userCredential *entity2.UserCredentials) (string, error) {
	username, err := s.repo.CheckUser(ctx, userCredential)
	if err != nil {
		return "", fmt.Errorf("-> s.Repo.CheckUser%v", err)
	}
	return username, nil
}

func (s *ShopService) CreateUser(ctx context.Context, userCredential *entity2.UserCredentials) (string, error) {
	username, err := s.repo.PutUser(ctx, userCredential)
	if err != nil {
		return "", fmt.Errorf("-> s.Repo.PutUser%v", err)
	}
	return username, nil
}
