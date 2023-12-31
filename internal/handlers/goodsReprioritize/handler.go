package goodsreprioritiize

import (
	"context"

	"github.com/catinapoke/go-microservice-example/internal/domain"
	"github.com/catinapoke/go-microservice-example/utils/serviceerrors"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	Priorities []domain.ItemPriority `json:"priorities"`
}

type Request struct {
	Id        int `query:"id"`
	ProjectId int `query:"projectId"`
	Priority  int `json:"newPriority"` // TODO: check args passing
}

func (r Request) Validate() error {
	if r.Priority < 1 {
		return serviceerrors.ErrValidation
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	priorities, err := h.Model.UpdatePriority(ctx, r.Id, r.ProjectId, r.Priority)
	return Response{Priorities: priorities}, err
}
