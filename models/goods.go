package models

import "time"

type Goods struct {
	Id          int64     `json:"id" db:"id"`
	ProjectId   int64     `json:"project_id" db:"project_id" binding:"required"`
	Name        string    `json:"name" db:"name" binding:"required"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed" binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type GoodsInput struct {
	Name        *string `json:"name" db:"name" binding:"required"`
	Description *string `json:"description" db:"description"`
}

type GoodsToLog struct {
	Id         int64 `json:"id"`
	CampaignId int64 `json:"campaignId"`
	Removed    bool  `json:"removed"`
}

type GoodsPriority struct {
	Id       int64 `json:"id" db:"id"`
	Priority int   `json:"priority" db:"priority"`
}
