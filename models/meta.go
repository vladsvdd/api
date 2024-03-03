package models

type Meta struct {
	Total   int `json:"total" db:"total"`
	Removed int `json:"removed" db:"removed"`
	Limit   int `json:"limit" db:"limit"`
	Offset  int `json:"offset" db:"offset"`
}
