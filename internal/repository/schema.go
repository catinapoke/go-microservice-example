package repository

import (
	"time"
)

type GoodsItem struct {
	Id          int       `db:"id"`
	ProjectId   int       `db:"projectid"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

type ProjectItem struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type GoodsPriority struct {
	Id        int `db:"id"`
	ProjectId int `db:"projectid"`
	Priority  int `db:"priority"`
}
