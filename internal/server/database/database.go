package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgerrcode"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
	"github.com/jackc/pgx"
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
		if pgerrcode.UniqueViolation == pgerrcode.ConnectionDoesNotExist {
			return true
		}
		if pgerrcode.UniqueViolation == pgerrcode.ConnectionException {
			return true
		}
	}
	return false
}

func retryQuery(ctx context.Context, query func() error) error {
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
	ctx.Done()
	return errors.New("failed after max retries")
}

func DBConnect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = retryQuery(context.Background(), func() error {
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
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DBSaveMetrics(db *sql.DB, metrics *storage.Storage) error {

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
