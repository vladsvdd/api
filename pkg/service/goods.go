package service

import (
	"api/models"
	"api/pkg/natsutil"
	"api/pkg/repository"
	"encoding/json"
	"fmt"
	"time"
)

const (
	GoodsNatsKey                = "goods.created"
	GoodsUpdatedPriorityNatsKey = "goods.updated.priority"
)

type GoodsService struct {
	repo        repository.Goods
	redisClient *repository.RedisClient
	natsClient  *natsutil.NATSClient
}

func NewGoodsService(
	repo repository.Goods,
	redisClient *repository.RedisClient,
	natsClient *natsutil.NATSClient,
) *GoodsService {
	return &GoodsService{
		repo:        repo,
		redisClient: redisClient,
		natsClient:  natsClient,
	}
}

func (s *GoodsService) Create(projectId int64, good models.GoodsInput) (models.Goods, error) {
	goodOutput, err := s.repo.Create(projectId, good)
	if err != nil {
		return goodOutput, err
	}

	if err := s.natsClient.Publish(GoodsNatsKey, goodOutput); err != nil {
		return models.Goods{}, fmt.Errorf("Ошибка публикации в nats, метод Create" + err.Error())
	}

	return goodOutput, nil
}

func (s *GoodsService) Update(id, projectId int64, goodInput models.GoodsInput) (models.Goods, error) {
	goodOutput, err := s.repo.Update(id, projectId, goodInput)
	if err != nil {
		return goodOutput, err
	}

	// Инвалидируем данные в redis
	marshalGood, err := json.Marshal(goodOutput)
	if err != nil {
		return goodOutput, err
	}

	if err := s.redisClient.Set(fmt.Sprintf("goods_%d_%d_Update", id, projectId), string(marshalGood), time.Minute); err != nil {
		return goodOutput, err
	}

	if err := s.natsClient.Publish(GoodsNatsKey, goodOutput); err != nil {
		return models.Goods{}, fmt.Errorf("Ошибка публикации в nats, метод Update" + err.Error())
	}

	return goodOutput, nil
}

func (s *GoodsService) UpdatePriority(id, projectId, newPriority int64) ([]models.GoodsPriority, error) {
	goodsOutput, err := s.repo.UpdatePriority(id, projectId, newPriority)
	if err != nil {
		return []models.GoodsPriority{}, err
	}

	// Инвалидируем данные в redis
	marshalGood, err := json.Marshal(goodsOutput)
	if err != nil {
		return []models.GoodsPriority{}, err
	}

	if err := s.redisClient.Set(fmt.Sprintf("goods_%d_%d_UpdatePriority", id, projectId), string(marshalGood), time.Minute); err != nil {
		return []models.GoodsPriority{}, err
	}

	if err := s.natsClient.Publish(GoodsUpdatedPriorityNatsKey, goodsOutput); err != nil {
		return []models.GoodsPriority{}, fmt.Errorf("Ошибка публикации в nats, метод Update" + err.Error())
	}

	return goodsOutput, nil
}

func (s *GoodsService) Delete(id, projectId int64) (models.GoodsToLog, error) {
	goodOutput, err := s.repo.Delete(id, projectId)
	if err != nil {
		return models.GoodsToLog{}, err
	}

	goodToLog := models.GoodsToLog{
		Id:         goodOutput.Id,
		CampaignId: goodOutput.ProjectId,
		Removed:    goodOutput.Removed,
	}

	// Инвалидируем данные в redis
	marshalGood, err := json.Marshal(goodToLog)
	if err != nil {
		return goodToLog, err
	}

	if err := s.redisClient.Set(fmt.Sprintf("goods_%d_%d_Delete", id, projectId), string(marshalGood), time.Minute); err != nil {
		return models.GoodsToLog{}, err
	}

	if err := s.natsClient.Publish(GoodsNatsKey, goodOutput); err != nil {
		return models.GoodsToLog{}, fmt.Errorf("Ошибка публикации в nats, метод Delete" + err.Error())
	}

	return goodToLog, nil
}

func (s *GoodsService) GetList(limit, offset int) ([]models.Goods, models.Meta, error) {
	keyGoods := fmt.Sprintf("goods_%d_%d_GetList", limit, offset)
	keyMeta := fmt.Sprintf("meta_%d_%d_GetList", limit, offset)

	// Получение данных из redis
	goods := []models.Goods{}
	goodsJson, err := s.redisClient.Get(keyGoods)
	if err == nil {
		_ = json.Unmarshal([]byte(goodsJson), &goods)
	}

	meta := models.Meta{}
	metaJson, err := s.redisClient.Get(keyMeta)
	if err == nil {
		_ = json.Unmarshal([]byte(metaJson), &meta)
	}

	// Возвращаем данные из redis
	if len(goods) > 0 && meta.Total != 0 {
		return goods, meta, nil
	}

	// Получение данных из БД, если в redis пусто
	goods, meta, err = s.repo.GetList(limit, offset)
	if err != nil {
		return []models.Goods{}, models.Meta{}, err
	}

	//Запись данных в redis
	marshalGoods, err := json.Marshal(goods)
	if err != nil {
		return []models.Goods{}, models.Meta{}, err
	}

	if err := s.redisClient.Set(keyGoods, string(marshalGoods), time.Minute); err != nil {
		return []models.Goods{}, models.Meta{}, err
	}

	marshalMeta, err := json.Marshal(meta)
	if err != nil {
		return []models.Goods{}, models.Meta{}, err
	}

	if err := s.redisClient.Set(keyMeta, string(marshalMeta), time.Minute); err != nil {
		return []models.Goods{}, models.Meta{}, err
	}

	return goods, meta, nil
}
