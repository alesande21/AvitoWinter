package service

type ShopRepo interface {
	Ping() error
}

type ShopService struct {
	Repo ShopRepo
}

func NewShopService(repo ShopRepo) *ShopService {
	return &ShopService{Repo: repo}
}
