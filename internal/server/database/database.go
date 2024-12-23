package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

func DBCreateTables(ctx context.Context, db *sql.DB) error {
	exec := `CREATE TABLE IF NOT EXISTS counter_metrics (
		name text NOT NULL UNIQUE,
		value bigint NOT NULL
		);

		CREATE TABLE IF NOT EXISTS gauge_metrics (
		name text NOT NULL UNIQUE,
		value double precision NOT NULL
		);`
	_, err := db.Exec(exec)
	return err
}

func DBSaveMetrics(ctx context.Context, db *sql.DB, metrics *storage.Storage) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Println("Rollback error:", rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()
	for k, v := range metrics.Gauge {

		_, err = tx.Exec("INSERT INTO gauge_metrics VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", k, v)
		if err != nil {
			return fmt.Errorf("error inserting gauge metric: %w", err)
		}
	}

	for k, v := range metrics.Counter {
		_, err := tx.Exec("INSERT INTO counter_metrics VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", k, v)
		if err != nil {
			return fmt.Errorf("error inserting counter metric: %w", err)
		}
	}
	return nil
}
