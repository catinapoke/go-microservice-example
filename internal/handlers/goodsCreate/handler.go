package goodscreate

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
	ProjectId int    `query:"projectId"`
	Name      string `json:"name"`
}

func (r Request) Validate() error {
	if r.Name == "" {
		return serviceerrors.ErrValidation
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	item, err := h.Model.Create(ctx, r.ProjectId, r.Name)
	return Response{item}, err
}
