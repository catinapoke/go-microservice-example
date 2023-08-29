package natslogs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/catinapoke/go-microservice-example/internal/domain"
	"github.com/catinapoke/go-microservice-example/internal/repository"
	"github.com/catinapoke/go-microservice-example/utils/logger"
	"github.com/nats-io/nats.go"
)

type Repository struct {
	conn    *nats.Conn
	next    domain.GoodsRepository
	subject string
}

type LogMessage struct {
	Id          int       `json:"Id"`
	ProjectId   int       `json:"ProjectId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Priority    int       `json:"Priority"`
	Removed     bool      `json:"Removed"`
	EventTime   time.Time `json:"EventTime"`
}

func New(conn *nats.Conn, repo domain.GoodsRepository, subject string) *Repository {
	return &Repository{
		conn:    conn,
		next:    repo,
		subject: subject,
	}
}

func (r *Repository) CreateItem(ctx context.Context, projectId int, name string) (*repository.GoodsItem, error) {
	item, err := r.next.CreateItem(ctx, projectId, name)

	if err == nil {
		r.tryPublish(item)
	}

	return item, err
}

func (r *Repository) GetItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	return r.next.GetItem(ctx, id, projectId)
}

func (r *Repository) UpdateItem(ctx context.Context, id int, projectId int, name string, description string) (*repository.GoodsItem, error) {
	item, err := r.next.UpdateItem(ctx, id, projectId, name, description)

	if err == nil {
		r.tryPublish(item)
	}

	return item, err
}

func (r *Repository) DeleteItem(ctx context.Context, id int, projectId int) (*repository.GoodsItem, error) {
	item, err := r.next.DeleteItem(ctx, id, projectId)
	if err == nil {
		r.tryPublish(item)
	}

	return item, err
}

func (r *Repository) ListItems(ctx context.Context, limit int, offset int) ([]repository.GoodsItem, error) {
	return r.next.ListItems(ctx, limit, offset)
}

func (r *Repository) Reprioritize(ctx context.Context, id int, projectId int, startPriority int) ([]repository.GoodsItem, error) {
	items, err := r.next.Reprioritize(ctx, id, projectId, startPriority)
	if err == nil {
		for _, item := range items {
			r.tryPublish(&item)
		}
	}

	return items, err
}

func (r *Repository) tryPublish(item *repository.GoodsItem) {
	if err := r.publish(item); err != nil {
		logger.Log.Warnf("nats publish err: %w", err)
	}
}

func (r *Repository) publish(item *repository.GoodsItem) error {
	logMsg := mapItemToLog(item)
	logger.Log.Infof("Preparing for publish: %+v", logMsg)

	data, err := json.Marshal(logMsg)

	if err != nil {
		return err
	}

	err = r.conn.Publish(r.subject, data)

	return err
}

func mapItemToLog(item *repository.GoodsItem) LogMessage {
	return LogMessage{
		Id:          item.Id,
		ProjectId:   item.ProjectId,
		Name:        item.Name,
		Description: item.Description,
		Priority:    item.Priority,
		Removed:     item.Removed,
		EventTime:   time.Now(),
	}
}
