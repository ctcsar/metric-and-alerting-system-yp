package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pressly/goose"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

const (
	maxRetries = 3
)

var retryDelays = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	5 * time.Second,
}

func isRetriableError(err error) bool {
	var pgErr *pgx.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.ConnectionDoesNotExist {
			return true
		}
		if pgErr.Code == pgerrcode.ConnectionException {
			return true
		}
	}
	return false
}

func retryQuery(query func() error) error {
	for i := 0; i < maxRetries; i++ {
		err := query()
		if err == nil {
			return nil
		}
		if !isRetriableError(err) {
			return err
		}
		time.Sleep(retryDelays[i])
	}
	return errors.New("failed after max retries")
}

func DBConnect(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DBMigrate(ctx context.Context, db *sql.DB) error {
	err := goose.Up(db, "../../migrations")
	if err != nil {
		return err
	}
	return nil
}

func DBSaveMetrics(ctx context.Context, db *sql.DB, metrics *storage.Storage) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
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
		err = retryQuery(func() error {
			_, err = tx.Exec("INSERT INTO gauge_metrics VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", k, v)
			return err
		})
		if err != nil {
			log.Println("Error inserting gauge metric:", err)
			continue
		}
	}

	for k, v := range metrics.Counter {
		err = retryQuery(func() error {
			_, err := tx.Exec("INSERT INTO counter_metrics VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", k, v)
			return err
		})
		if err != nil {
			log.Println("Error inserting counter metric:", err)
			continue
		}
	}
	if err != nil {
		return fmt.Errorf("error saving metrics: %w", err)
	}
	return nil
}
