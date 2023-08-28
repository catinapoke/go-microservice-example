-- +goose Up
-- +goose StatementBegin
CREATE TABLE Projects (
  	id serial,
	name text,
	created_at timestamp,
  	primary key (id)
);

CREATE TABLE Goods (
  	id integer,
  	projectId integer,
	name text,
  	description text,
  	priority integer,
  	removed boolean,
	created_at timestamp,
  	primary key (id, projectId)
);

CREATE INDEX ON Goods (id, projectId, name);

INSERT INTO Projects VALUES (1, 'Первая запись', NOW())
ON CONFLICT DO NOTHING;
-- +goose StatementEnd
-- 

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Projects;
DROP TABLE IF EXISTS Goods;
-- +goose StatementEnd