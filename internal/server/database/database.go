package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
)

// const (
// 	maxRetries = 3
// )

// var retryDelays = []time.Duration{
// 	1 * time.Second,
// 	3 * time.Second,
// 	5 * time.Second,
// }

// func isRetriableError(err error) bool {
// 	var pgErr *pgx.PgError
// 	if errors.As(err, &pgErr) {
// 		if pgErr.Code == pgerrcode.ConnectionDoesNotExist {
// 			return true
// 		}
// 		if pgErr.Code == pgerrcode.ConnectionException {
// 			return true
// 		}
// 	}
// 	return false
// }

// func retryQuery(query func() error) error {
// 	for i := 0; i < maxRetries; i++ {
// 		err := query()
// 		if err == nil {
// 			return nil
// 		}
// 		if !isRetriableError(err) {
// 			return err
// 		}
// 		time.Sleep(retryDelays[i])
// 	}
// 	return errors.New("failed after max retries")
// }

func DBConnect(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func DBMigrate(ctx context.Context, db *sql.DB) error {
	// err := goose.Up(db, "../../migrations")
	// if err != nil {
	// 	return err
	// }
	// return nil
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

func insertMetrics(ctx context.Context, tx *sql.Tx, tableName string, metrics map[string]interface{}) error {
	placeholders := []string{}
	args := []interface{}{}
	argIndex := 1

	for k, v := range metrics {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1))
		args = append(args, k, v)
		argIndex += 2
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (name, value) VALUES %s ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value",
		tableName,
		strings.Join(placeholders, ", "),
	)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert metrics into %s: %w", tableName, err)
	}
	return nil
}
func DBSaveMetrics(ctx context.Context, db *sql.DB, metrics *storage.Storage) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context expired before transaction start: %w", ctx.Err())
	}
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

	metrics.Mutex.RLock()
	defer metrics.Mutex.RUnlock()

	gaugeMetrics := map[string]interface{}{}
	for k, v := range metrics.Gauge {
		gaugeMetrics[k] = v
	}

	counterMetrics := map[string]interface{}{}
	for k, v := range metrics.Counter {
		counterMetrics[k] = v
	}

	err = insertMetrics(ctx, tx, "gauge_metrics", gaugeMetrics)
	if err != nil {
		return err
	}

	err = insertMetrics(ctx, tx, "counter_metrics", counterMetrics)
	if err != nil {
		return err
	}

	return nil
}
