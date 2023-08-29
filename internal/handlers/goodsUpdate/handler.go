package goodsupdate

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/domain"
	"github.com/catinapoke/go-microservice-example/utils/serviceerrors"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	domain.Item
}

type Request struct {
	Id          int    `query:"id"`
	ProjectId   int    `query:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r Request) Validate() error {
	if r.Name == "" {
		return serviceerrors.ErrValidation
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	item, err := h.Model.Update(ctx, r.Id, r.ProjectId, r.Name, r.Description)
	return Response{item}, err
}
