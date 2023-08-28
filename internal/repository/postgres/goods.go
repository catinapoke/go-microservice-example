package postgres

import (
	"context"
	"fmt"

	"github.com/catinapoke/go-microservice-example/internal/repository"
	"github.com/catinapoke/go-microservice-example/utils/tx"

	"github.com/georgysavva/scany/pgxscan"
)

type Repository struct {
	provider tx.DBProvider
}

const (
	GoodsTableName    = "Goods"
	ProjectsTableName = "Projects"
)

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

func (r *Repository) CreateItem(ctx context.Context, projectId int, name string) (*repository.GoodsItem, error) {
	query := `
	with items AS(
		SELECT id, projectId, name, priority
			FROM goods
			WHERE projectId = $1
		   
			union all
		
		SELECT 0 as id, $1 as projectId, $2 as name, 0 as priority
		)
		insert INTO goods(id, projectId, name, priority)
		select max(id) + 1, $1, $2, max(priority) + 1  from items
		returning *;
	`

	db := r.provider.GetDB(ctx)
	row := db.QueryRow(ctx, query, projectId, name)

	var result repository.GoodsItem
	err := row.Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("create: query row scan: %w", err)
	}

	return &result, nil
}

func (r *Repository) GetItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	query := `
	select *
	from goods
	where id = $1 AND projectId = $2;
	`

	db := r.provider.GetDB(ctx)
	row := db.QueryRow(ctx, query, id, projectId)

	var result repository.GoodsItem
	err := row.Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("get: query row scan: %w", err)
	}

	return &result, nil
}

func (r *Repository) UpdateItem(ctx context.Context, id int, projectId int, name string, description string) (*repository.GoodsItem, error) {
	query := `
	update goods
	set name = $3, description = $4
	where id = $1 AND projectId = $2
	RETURNING *;
	`

	db := r.provider.GetDB(ctx)
	row := db.QueryRow(ctx, query, id, projectId, name, description)

	var result repository.GoodsItem
	err := row.Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("update: query row scan: %w", err)
	}

	return &result, nil
}

func (r *Repository) DeleteItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	query := `
	update goods
	set removed = true
	where id = $1 AND projectId = $2
	RETURNING *;
	`

	db := r.provider.GetDB(ctx)
	row := db.QueryRow(ctx, query, id, projectId)

	var result repository.GoodsItem
	err := row.Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("delete: query row scan: %w", err)
	}

	return &result, nil
}

func (r *Repository) ListItems(ctx context.Context, limit int, offset int) ([]repository.GoodsItem, error) {
	query := fmt.Sprintf(`
		SELECT * 
		FROM Goods
		WHERE id >= %d
		ORDER BY projectId, id ASC
		LIMIT %d;
	`, offset, limit)

	db := r.provider.GetDB(ctx)
	rows, err := db.Query(ctx, query)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("list: sql query: %w", err)
	}

	var result []repository.GoodsItem = make([]repository.GoodsItem, 0)

	for rows.Next() {
		var item repository.GoodsItem

		if err := pgxscan.ScanRow(&item, rows); err != nil {
			return []repository.GoodsItem{}, fmt.Errorf("list: query scan: %w", err)
		}

		result = append(result, item)
	}

	return result, nil
}

func (r *Repository) Reprioritize(ctx context.Context, id int, projectId int, startPriority int) ([]repository.GoodsPriority, error) {
	query := `
	update goods
	set priority = id - $1 + $3
	where id >= $1 AND projectId = $2 AND priority != id - $1 + $3
	RETURNING id, projectId, priority;
	`

	db := r.provider.GetDB(ctx)
	rows, err := db.Query(ctx, query, id, projectId, startPriority)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("reprioritize: sql query: %w", err)
	}

	var result []repository.GoodsPriority = make([]repository.GoodsPriority, 0)

	for rows.Next() {
		var priority repository.GoodsPriority

		if err := pgxscan.ScanRow(&priority, rows); err != nil {
			return []repository.GoodsPriority{}, fmt.Errorf("reprioritize: query scan: %w", err)
		}

		result = append(result, priority)
	}

	return result, nil
}
