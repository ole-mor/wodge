package postgres

import (
	"context"
	"database/sql"
	"wodge/internal/services"

	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	db *sql.DB
}

func NewPostgresDriver(dsn string) (*PostgresDriver, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDriver{db: db}, nil
}

// Ensure PostgresDriver implements services.DatabaseService
var _ services.DatabaseService = (*PostgresDriver)(nil)

func (p *PostgresDriver) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold values for each column
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		// Create map for this row
		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			// Handle byte slices (often returned by driver for strings)
			if b, ok := val.([]byte); ok {
				rowMap[colName] = string(b)
			} else {
				rowMap[colName] = val
			}
		}
		results = append(results, rowMap)
	}

	return results, nil
}

func (p *PostgresDriver) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
