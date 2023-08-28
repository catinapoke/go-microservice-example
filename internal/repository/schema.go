package repository

import (
	"time"
)

type GoodsItem struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProjectItem struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type GoodsPriority struct {
	Id        int `json:"id"`
	ProjectId int `json:"projectId"`
	Priority  int `json:"priority"`
}
