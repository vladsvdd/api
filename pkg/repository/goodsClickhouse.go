package repository

import (
	"api/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type GoodsClickhouse struct {
	dbClickhouse *sqlx.DB
}

func NewGoodsClickhouse(dbClickhouse *sqlx.DB) *GoodsClickhouse {
	return &GoodsClickhouse{
		dbClickhouse: dbClickhouse,
	}
}

func (p *GoodsClickhouse) Create(good models.Goods) error {
	tx, err := p.dbClickhouse.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query := fmt.Sprintf(`
        INSERT INTO %s 
            (id, 
             ProjectId, 
             Name, 
             Description, 
             Priority, 
             Removed, 
             EventTime) 
        VALUES 
            (?,?,?,?,?,?,?) 
        `, goodsTable)

	// Выполните запрос с помощью QueryRow, чтобы получить вставленную запись.
	_, err = tx.Exec(query,
		good.Id,
		good.ProjectId,
		good.Name,
		good.Description,
		good.Priority,
		good.Removed,
		good.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert new record: %v", err)
	}

	// Если все успешно, завершите транзакцию.
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (p *GoodsClickhouse) CreateBatch(goods []models.Goods) error {
	tx, err := p.dbClickhouse.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	if err = p.createBatchRecord(tx, goods); err != nil {
		return err
	}

	return nil
}

func (p *GoodsClickhouse) createBatchRecord(tx *sqlx.Tx, records []models.Goods) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, 
			ProjectId, 
			Name, 
			Description, 
			Priority, 
			Removed, 
			EventTime
		) VALUES (
			:id, 
			:ProjectId, 
			:Name, 
			:Description, 
			:Priority, 
			:Removed, 
			:EventTime
		)
	`, goodsTable)

	// Подготовка вне цикла
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return errors.New("Ошибка при подготовке запроса: " + err.Error())
	}
	defer stmt.Close()

	for _, product := range records {
		removedValue := uint8(0)
		if product.Removed {
			removedValue = 1
		}

		good := map[string]interface{}{
			"id":          product.Id,
			"ProjectId":   product.ProjectId,
			"Name":        product.Name,
			"Description": product.Description,
			"Priority":    product.Priority,
			"Removed":     removedValue,
			"EventTime":   product.CreatedAt,
		}

		_, err := stmt.Exec(good)
		if err != nil {
			return errors.New("Ошибка добавления пакета " + err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Ошибка при подтверждении транзакции:" + err.Error())
	}

	return nil
}

func (p *GoodsClickhouse) CreatePriorityBatch(goods []models.GoodsPriority) error {
	tx, err := p.dbClickhouse.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	if err = p.createPriorityBatch(tx, goods); err != nil {
		return err
	}

	return nil
}

func (p *GoodsClickhouse) createPriorityBatch(tx *sqlx.Tx, records []models.GoodsPriority) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			id, 
			Priority
		) VALUES (
			:id, 
			:Priority
		)
	`, goodsTable)

	// Подготовка вне цикла
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return errors.New("Ошибка при подготовке запроса: " + err.Error())
	}
	defer stmt.Close()

	for _, product := range records {

		good := map[string]interface{}{
			"id":       product.Id,
			"Priority": product.Priority,
		}

		_, err := stmt.Exec(good)
		if err != nil {
			return errors.New("Ошибка добавления пакета " + err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Ошибка при подтверждении транзакции:" + err.Error())
	}

	return nil
}
