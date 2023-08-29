package domain

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/repository"
)

func (m *Model) Create(ctx context.Context, projectId int, name string) (Item, error) {
	var item *repository.GoodsItem
	err := m.txm.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var err error
		item, err = m.repo.CreateItem(ctx, projectId, name)
		return err
	})

	if err != nil {
		return Item{}, err
	}

	return mapRepoItemToModel(*item), nil
}
