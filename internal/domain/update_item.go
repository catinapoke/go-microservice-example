package domain

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/repository"
)

func (m *Model) Update(ctx context.Context, id int, projectId int, name string, description string) (Item, error) {
	var item *repository.GoodsItem
	err := m.txm.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		var err error
		item, err = m.repo.UpdateItem(ctx, id, projectId, name, description)
		return err
	})

	if err != nil {
		return Item{}, err
	}

	return mapRepoItemToModel(*item), nil
}
