package repository

import (
	"api/models"
	"github.com/jmoiron/sqlx"
)

type Goods interface {
	Create(projectId int64, good models.GoodsInput) (models.Goods, error)
	Update(id, projectId int64, good models.GoodsInput) (models.Goods, error)
	UpdatePriority(id, projectId, newPriority int64) ([]models.GoodsPriority, error)
	Delete(id, projectId int64) (models.Goods, error)
	GetList(limit, offset int) ([]models.Goods, models.Meta, error)
}

type GoodsToClickhouse interface {
	Create(good models.Goods) error
	CreateBatch(goods []models.Goods) error
	CreatePriorityBatch(goods []models.GoodsPriority) error
}

type Repository struct {
	Goods
	GoodsToClickhouse
}

func NewRepository(db *sqlx.DB, dbClickhouse *sqlx.DB) *Repository {
	return &Repository{
		Goods:             NewGoodsPostgres(db),
		GoodsToClickhouse: NewGoodsClickhouse(dbClickhouse),
	}
}
