package service

import (
	"api/models"
	"api/pkg/natsutil"
	"api/pkg/repository"
)

type Goods interface {
	Create(projectId int64, good models.GoodsInput) (models.Goods, error)
	Update(id, projectId int64, good models.GoodsInput) (models.Goods, error)
	UpdatePriority(id, projectId, newPriority int64) ([]models.GoodsPriority, error)
	Delete(id, projectId int64) (models.GoodsToLog, error)
	GetList(limit, offset int) ([]models.Goods, models.Meta, error)
}

type Service struct {
	Goods
}

func NewService(
	repos *repository.Repository,
	redisClient *repository.RedisClient,
	natsClient *natsutil.NATSClient,
) *Service {
	return &Service{
		Goods: NewGoodsService(repos.Goods, redisClient, natsClient),
	}
}
