package goodsremove

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	Id        int  `json:"id"`
	ProjectId int  `json:"projectId"`
	Removed   bool `json:"removed"`
}

type Request struct {
	Id        int `query:"id"`
	ProjectId int `query:"projectId"`
}

func (r Request) Validate() error {
	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	item, err := h.Model.Remove(ctx, r.Id, r.ProjectId)

	if err != nil {
		return Response{}, err
	}

	return Response{
		Id:        item.Id,
		ProjectId: item.ProjectId,
		Removed:   item.Removed,
	}, nil
}
