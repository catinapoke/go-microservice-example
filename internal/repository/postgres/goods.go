package postgres

import (
	"context"
	"fmt"

	"github.com/catinapoke/go-microservice-example/internal/repository"
	"github.com/catinapoke/go-microservice-example/utils/serviceerrors"
	"github.com/catinapoke/go-microservice-example/utils/tx"
	"github.com/jackc/pgx/v4"

	"github.com/georgysavva/scany/pgxscan"
)

type Repository struct {
	provider tx.DBProvider
}

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
		returning id, projectId, name, description, priority, removed, created_at;
	`

	db := r.provider.GetDB(ctx)

	if err := r.lockProjectItems(ctx, projectId, db); err != nil {
		return nil, err
	}

	row := db.QueryRow(ctx, query, projectId, name)

	var result repository.GoodsItem
	err := row.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Description, &result.Priority, &result.Removed, &result.CreatedAt)
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
	err := row.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Description, &result.Priority, &result.Removed, &result.CreatedAt)
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
	RETURNING id, projectId, name, description, priority, removed, created_at;
	`

	db := r.provider.GetDB(ctx)

	if err := r.lockItem(ctx, id, projectId, db); err != nil {
		return nil, err
	}

	row := db.QueryRow(ctx, query, id, projectId, name, description)

	var result repository.GoodsItem
	err := row.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Description, &result.Priority, &result.Removed, &result.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, serviceerrors.ErrNotFound
		}

		return nil, fmt.Errorf("update: query row scan: %w", err)
	}

	return &result, nil
}

func (r *Repository) DeleteItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	query := `
	update goods
	set removed = true
	where id = $1 AND projectId = $2
	RETURNING id, projectId, name, description, priority, removed, created_at;
	`

	db := r.provider.GetDB(ctx)

	if err := r.lockItem(ctx, id, projectId, db); err != nil {
		return nil, err
	}

	row := db.QueryRow(ctx, query, id, projectId)

	var result repository.GoodsItem
	err := row.Scan(&result.Id, &result.ProjectId, &result.Name, &result.Description, &result.Priority, &result.Removed, &result.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, serviceerrors.ErrNotFound
		}

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

func (r *Repository) Reprioritize(ctx context.Context, id int, projectId int, startPriority int) ([]repository.GoodsItem, error) {
	query := `
	update goods
	set priority = id - $1 + $3
	where id >= $1 AND projectId = $2 AND priority != id - $1 + $3
	RETURNING id, projectId, name, description, priority, removed, created_at;
	`

	db := r.provider.GetDB(ctx)

	if err := r.lockProjectItems(ctx, projectId, db); err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, query, id, projectId, startPriority)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("reprioritize: sql query: %w", serviceerrors.ErrInternal)
	}

	var result []repository.GoodsItem = make([]repository.GoodsItem, 0)

	for rows.Next() {
		var priority repository.GoodsItem

		if err := pgxscan.ScanRow(&priority, rows); err != nil {
			return []repository.GoodsItem{}, fmt.Errorf("reprioritize: query scan: %w", serviceerrors.ErrInternal)
		}

		result = append(result, priority)
	}

	return result, nil
}

func (r *Repository) lockItem(ctx context.Context, id int, projectId int, db tx.Querier) error {
	query := `
		select *
		from goods
		where id = $1 AND projectId = $2
		for update;
	`

	_, err := db.Exec(ctx, query, id, projectId)
	if err != nil {
		return fmt.Errorf("exec lock item: %w", serviceerrors.ErrInternal)
	}

	return nil
}

func (r *Repository) lockProjectItems(ctx context.Context, projectId int, db tx.Querier) error {
	query := `
	select *
	from goods
	where projectId = $1
	for update;
`

	count, err := db.Exec(ctx, query, projectId)

	if err != nil {
		return fmt.Errorf("exec lock project items: %w", serviceerrors.ErrInternal)
	}

	if count.RowsAffected() == 0 {
		return serviceerrors.ErrNotFound
	}

	return nil
}
