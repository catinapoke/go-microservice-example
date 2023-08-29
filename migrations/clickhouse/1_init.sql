-- +goose Up
-- +goose StatementBegin
CREATE TABLE Goods (
  	Id Int64,
  	ProjectId Int64,
	Name String,
  	Description String,
  	Priority Int64,
  	Removed Boolean,
	EventTime DateTime
) ENGINE = NATS
	SETTINGS nats_url = 'nats:4222',
		nats_subjects = 'GoodsLogs',
		nats_format = 'JSON';
-- +goose StatementEnd
-- 

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Goods;
-- +goose StatementEnd