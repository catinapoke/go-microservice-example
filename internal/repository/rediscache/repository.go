package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/catinapoke/go-microservice-example/internal/domain"
	"github.com/catinapoke/go-microservice-example/internal/repository"
	"github.com/catinapoke/go-microservice-example/utils/logger"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
	next   domain.GoodsRepository
	key    string
}

func New(client *redis.Client, repo domain.GoodsRepository, key string) *Repository {
	return &Repository{
		client: client,
		next:   repo,
		key:    key,
	}
}

func (r *Repository) CreateItem(ctx context.Context, projectId int, name string) (*repository.GoodsItem, error) {
	r.client.Del(ctx, r.key)
	return r.next.CreateItem(ctx, projectId, name)
}

func (r *Repository) GetItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	r.client.Del(ctx, r.key)
	return r.next.GetItem(ctx, id, projectId)
}

func (r *Repository) UpdateItem(ctx context.Context, id int, projectId int, name string, description string) (*repository.GoodsItem, error) {
	r.client.Del(ctx, r.key)
	return r.next.UpdateItem(ctx, id, projectId, name, description)
}

func (r *Repository) DeleteItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	r.client.Del(ctx, r.key)
	return r.next.DeleteItem(ctx, id, projectId)
}

func (r *Repository) ListItems(ctx context.Context, limit int, offset int) ([]repository.GoodsItem, error) {
	var items []repository.GoodsItem

	items, err := r.getList(ctx, limit, offset)

	if err != nil {
		logger.Log.Infof("Missed cache for /list limit:%d, offset:%d err: %v", limit, offset, err)
		items, err = r.next.ListItems(ctx, limit, offset)

		if err != nil {
			return nil, err
		}

		r.setList(ctx, limit, offset, items)
	}

	return items, nil
}

func (r *Repository) Reprioritize(ctx context.Context, id int, projectId int, startPriority int) ([]repository.GoodsItem, error) {
	r.client.Del(ctx, r.key)
	return r.next.Reprioritize(ctx, id, projectId, startPriority)
}

func (r *Repository) setList(ctx context.Context, limit int, offset int, items []repository.GoodsItem) error {
	data, err := json.Marshal(items)

	if err != nil {
		return err
	}

	res, err := r.client.HSet(ctx, r.key, r.getKeyField(limit, offset), data).Result()
	logger.Log.Infof("Cache set request: %v %v", res, err)

	if err != nil {
		r.client.Expire(ctx, r.key, time.Minute)
	}

	return err
}

func (r *Repository) getList(ctx context.Context, limit int, offset int) ([]repository.GoodsItem, error) {
	var items []repository.GoodsItem

	data, err := r.client.HGet(ctx, r.key, r.getKeyField(limit, offset)).Result()

	if err != nil {
		logger.Log.Infof("getList: missing cache: %w", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &items)

	if err != nil {
		logger.Log.Infof("getList: unmarshal: %w", err)
		return nil, err
	}

	return items, nil
}

func (r *Repository) getKeyField(limit int, offset int) string {
	return fmt.Sprintf("%d - %d", limit, offset)
}
