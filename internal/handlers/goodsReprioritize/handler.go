package goodsreprioritiize

import (
	"context"
	"errors"
	"log"

	"github.com/catinapoke/go-microservice-example/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	Priorities []domain.ItemPriority `json:"priorities"`
}

type Request struct {
	Id        int `path:"id"`
	ProjectId int `path:"projectId"`
	Priority  int `query:"newPriority"` // TODO: check args passing
}

var (
	ErrWrongInput = errors.New("wrong input")
)

func (r Request) Validate() error {
	if r.Priority < 1 {
		return ErrWrongInput
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	log.Printf("%+v", r)

	priorities, err := h.Model.UpdatePriority(ctx, r.Id, r.ProjectId, r.Priority)
	return Response{Priorities: priorities}, err
}
