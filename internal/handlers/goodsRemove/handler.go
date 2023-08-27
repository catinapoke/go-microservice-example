package goodsremove

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
	Id        int  `json:"id"`
	ProjectId int  `json:"projectId"` // TODO: Check if it's really campaignId
	Removed   bool `json:"removed"`   // true
}

type Request struct {
	Id        int `path:"id"`
	ProjectId int `path:"projectId"`
}

var (
	ErrWrongInput = errors.New("wrong input")
)

func (r Request) Validate() error {
	return nil
}

func (h *Handler) Handle(ctx context.Context, r Request) (Response, error) {
	log.Printf("%+v", r)

	err := h.Model.Remove(ctx, r.Id, r.ProjectId)
	return Response{
		Id:        r.Id,
		ProjectId: r.ProjectId,
		Removed:   true,
	}, err
}
