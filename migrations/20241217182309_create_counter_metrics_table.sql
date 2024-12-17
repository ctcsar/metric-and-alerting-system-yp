-- +goose Up
CREATE TABLE counter_metrics (
	         name text NOT NULL UNIQUE,
	         value bigint NOT NULL
	         );
             
-- +goose Down
DROP TABLE counter_metrics;
