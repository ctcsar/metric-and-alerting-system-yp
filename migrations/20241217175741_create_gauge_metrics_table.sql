-- +goose Up
CREATE TABLE gauge_metrics (
            name text NOT NULL UNIQUE,
            value double precision NOT NULL
            );

-- +goose Down
DROP TABLE gauge_metrics;