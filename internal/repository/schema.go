package repository

import (
	"time"
)

type GoodsItem struct {
	Id          int       `db:"id"`
	ProjectId   int       `db:"projectId"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"createdAt"`
}

type ProjectItem struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"createdAt"`
}

type GoodsPriority struct {
	Id        int `db:"id"`
	ProjectId int `db:"projectId"`
	Priority  int `db:"priority"`
}
