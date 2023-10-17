package repository

import (
	"database/sql"
	"os"
	"statForMarket/internal/model"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type Repository struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{
		Conn: conn,
	}
}

func InitRepository() *sql.DB {
	return clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{os.Getenv("HOST") + ":" + os.Getenv("DB_PORT")},
		Auth: clickhouse.Auth{
			Database: os.Getenv("DATABASE"),
			Username: os.Getenv("USERNAME"),
			Password: os.Getenv("PASSWORD"),
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 30 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Protocol: clickhouse.HTTP,
	})
}

func (r *Repository) TestEvents(events []*model.Event) error {
	var values []interface{}
	query := `INSERT INTO events (eventID, eventType, userID, eventTime, payload) VALUES `
	for _, event := range events {
		query += "(?, ?, ?, ?, ?), "
		values = append(values, event.EventID, event.EventType, event.UserID, event.EventTime, event.Payload)
	}
	query = strings.TrimSuffix(query, ", ") + ";"

	_, err := r.Conn.Exec(query, values...)
	return err
}

func (r *Repository) CreateEvent(event *model.Event) error {
	query := `INSERT INTO events (eventID, eventType, userID, eventTime, payload) VALUES (?, ?, ?, ?, ?);`
	_, err := r.Conn.Exec(query, event.EventID, event.EventType, event.UserID, event.EventTime, event.Payload)
	return err
}
