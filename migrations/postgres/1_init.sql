-- +goose Up
-- +goose StatementBegin
CREATE TABLE Projects (
  	id serial,
	name text NOT NULL,
	created_at timestamp DEFAULT NOW(),
  	primary key (id)
);

CREATE TABLE Goods (
  	id integer,
  	projectId integer,
	name text NOT NULL,
  	description text DEFAULT '',
  	priority integer DEFAULT 1,
  	removed boolean DEFAULT false,
	created_at timestamp DEFAULT NOW(),
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