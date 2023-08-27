package domain

import "time"

type Item struct {
	Id          int       `json:"id"`
	ProjectId   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ItemPriority struct {
	Id       int `json:"id"`
	Priority int `json:"priority"`
}

type GoodsCreator interface {
	Create(projectId int, name string) (Item, error)
}

type GoodsUpdater interface {
	Update(id int, projectId int, name string, description string) (Item, error)
}

type GoodsRemover interface {
	Remove(id int, projectId int) error
}

type GoodsListing interface {
	List(limit int, offset int) ([]Item, error)
}

type GoodsPrioritizer interface {
	UpdatePriority(id int, projectId int, newPriority int) ([]ItemPriority, error)
}

type Model struct {
}

func New() *Model {
	return &Model{
	}
}
