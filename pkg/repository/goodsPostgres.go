package repository

import (
	"api/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type GoodsPostgres struct {
	db *sqlx.DB
}

func NewGoodsPostgres(db *sqlx.DB) *GoodsPostgres {
	return &GoodsPostgres{
		db: db,
	}
}

func (p *GoodsPostgres) Create(projectId int64, good models.GoodsInput) (models.Goods, error) {
	var insertedGood models.Goods
	tx, err := p.db.Beginx()
	if err != nil {
		return insertedGood, fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := fmt.Sprintf(`
        INSERT INTO %s 
            (project_id,
            name,
            description) 
        VALUES 
            ($1, $2, $3) 
        RETURNING *`, goodsTable)

	err = tx.QueryRowx(query, projectId, good.Name, good.Description).StructScan(&insertedGood)
	if err != nil {
		return insertedGood, fmt.Errorf("failed to insert new record: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return insertedGood, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return insertedGood, nil
}

func (p *GoodsPostgres) Update(id, projectId int64, good models.GoodsInput) (models.Goods, error) {
	// Начинаем транзакцию с блокировкой на чтение записи. Уровень изоляции Serializable
	tx, err := p.db.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return models.Goods{}, fmt.Errorf("не удалось запустить транзакцию: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	// Проверить валидность полей.
	if err := validateGoodsInput(good); err != nil {
		return models.Goods{}, fmt.Errorf("недопустимый ввод: %v", err)
	}

	queryExist := fmt.Sprintf(`
		SELECT 
		    id
		FROM 
		    %s
		WHERE
		    id = $1
			AND project_id = $2
	`, goodsTable)

	var exist int64
	if err := tx.QueryRowxContext(context.Background(), queryExist, id, projectId).Scan(&exist); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Goods{}, errors.New("errors.good.notFound")
		}
		return models.Goods{}, err
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET
		    name = COALESCE($1, name),
		    description = COALESCE($2, description) 
        WHERE
		    id = $3
		    AND project_id = $4
        RETURNING *`, goodsTable)

	var goodUpdated models.Goods
	err = tx.QueryRowxContext(context.Background(), query, good.Name, good.Description, id, projectId).StructScan(&goodUpdated)
	if err != nil {
		return models.Goods{}, fmt.Errorf("не удалось вставить новую запись: %v", err)
	}

	// Если все успешно, завершите транзакцию.
	err = tx.Commit()
	if err != nil {
		return models.Goods{}, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	return goodUpdated, nil
}

func (p *GoodsPostgres) UpdatePriority(id, projectId, newPriority int64) ([]models.GoodsPriority, error) {
	tx, err := p.db.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return []models.GoodsPriority{}, fmt.Errorf("не удалось запустить транзакцию: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	queryExist := fmt.Sprintf(`
		SELECT 
		    priority
		FROM 
		    %s
		WHERE
		    id = $1
			AND project_id = $2
	`, goodsTable)

	var priorityOld int64
	if err := tx.QueryRowxContext(context.Background(), queryExist, id, projectId).Scan(&priorityOld); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.GoodsPriority{}, errors.New("errors.good.notFound")
		}
		return []models.GoodsPriority{}, err
	}

	// Если уровень приоритета уже установле до нужного, то не зачем делать лишние запросы
	if priorityOld == newPriority {
		return []models.GoodsPriority{}, fmt.Errorf("уровень приоритета %d уже установлен", priorityOld)
	}

	good := models.GoodsPriority{}
	queryUpdateOne := fmt.Sprintf(`
		UPDATE 
			%s
		SET 
		    priority = $1
		WHERE
			id = $2
			AND project_id = $3
		RETURNING 
		    id, priority`, goodsTable)
	err = tx.QueryRowx(queryUpdateOne, newPriority, id, projectId).StructScan(&good)
	if err != nil {
		return []models.GoodsPriority{}, fmt.Errorf("не удалось обновить запись: %v", err)
	}

	goods := []models.GoodsPriority{}
	queryUpdateOther := fmt.Sprintf(`
		UPDATE 
		    %s
		SET 
		    priority = priority + 1
		WHERE 
		    priority >= $1
		    AND id != $2
        RETURNING 
		    id, priority`, goodsTable)

	err = tx.Select(&goods, queryUpdateOther, newPriority, id)
	if err != nil {
		return []models.GoodsPriority{}, fmt.Errorf("не удалось обновить записи: %v", err)
	}
	goods = append(goods, good)

	err = tx.Commit()
	if err != nil {
		return []models.GoodsPriority{}, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	return goods, nil
}

func validateGoodsInput(good models.GoodsInput) error {
	if good.Name == nil {
		return errors.New("поле Name не должно быть пустым")
	}
	return nil
}

func (p *GoodsPostgres) Delete(id, projectId int64) (models.Goods, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return models.Goods{}, err
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	// Выборка данных для последующего удаления
	var goodDeleted models.Goods
	err = tx.Get(&goodDeleted, fmt.Sprintf(`
		SELECT 
		    * 
		FROM 
		    %s 
		WHERE 
		    id = $1 
		    AND project_id = $2`, goodsTable),
		id, projectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Нет данных для удаления
			return models.Goods{}, errors.New("errors.good.notFound")
		}
		return models.Goods{}, err
	}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE
		    id = $1
			AND project_id = $2
	`, goodsTable)

	tx.QueryRow(query, id, projectId)

	err = tx.Commit()
	if err != nil {
		return models.Goods{}, fmt.Errorf("не удалось зафиксировать транзакцию: %v", err)
	}

	goodDeleted.Removed = true

	return goodDeleted, nil
}

func (p *GoodsPostgres) GetList(limit, offset int) ([]models.Goods, models.Meta, error) {
	meta := models.Meta{}
	goods := []models.Goods{}
	const (
		LIMIT     = 10
		OFFSET    = 1
		TOP_LIMIT = 100
	)

	if limit <= 0 {
		limit = LIMIT
	}

	if offset <= 0 {
		offset = OFFSET
	}

	//устанавливаем разумное ограничение на уровне БД в 100 записей за 1 раз
	if limit > TOP_LIMIT {
		limit = TOP_LIMIT
	}

	if err := p.db.Get(&meta, fmt.Sprintf(`
			SELECT
				COUNT(1) AS total,
				COUNT(CASE WHEN removed = true THEN 1 END) AS removed
			FROM
				%s;
		`, goodsTable)); err != nil {
		return goods, meta, err
	}

	meta.Limit = limit
	meta.Offset = offset

	// Постраничное чтение
	query := fmt.Sprintf(`
			SELECT
			    subquery.id,
				subquery.project_id,
				subquery.name,
				subquery.description,
				subquery.priority,
				subquery.removed,
				subquery.created_at
			FROM (
				SELECT
					*,
					ROW_NUMBER() OVER (ORDER BY id) as row_num
				FROM
					%s
			) AS subquery
			WHERE
				subquery.row_num > $1 AND subquery.row_num <= $2
			`, goodsTable)

	if err := p.db.Select(&goods, query, offset, offset+limit); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goods, meta, errors.New("errors.good.notFound")
		}
		return goods, meta, err
	}

	return goods, meta, nil
}
