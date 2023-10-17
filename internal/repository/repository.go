package repository

import (
	"database/sql"
	"fmt"
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

func (r *Repository) Events(filter *model.EventFilter) ([]*model.Event, error) {
	events := make([]*model.Event, 0)
	query := `SELECT * FROM events `
	if filter.EventType != "" && filter.From == "" {
		query += fmt.Sprintf("WHERE eventType = '%s' ", filter.EventType)
	}
	if filter.EventType == "" && filter.From != "" {
		query += fmt.Sprintf("WHERE eventTime > '%s' AND eventTime < '%s'", filter.From, filter.To)
	}
	if filter.EventType != "" && filter.From != "" {
		query += fmt.Sprintf("WHERE eventType = '%s' AND eventTime > '%s' AND eventTime < '%s'", filter.EventType, filter.From, filter.To)
	}
	fmt.Printf("query: %v\n", query)

	rows, err := r.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := new(model.Event)
		if err := rows.Scan(&event.EventID, &event.EventType, &event.UserID, &event.EventTime, &event.Payload); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, err
}

func (r *Repository) CreateEvent(event *model.Event) error {
	query := `INSERT INTO events (eventID, eventType, userID, eventTime, payload) VALUES (?, ?, ?, ?, ?);`
	_, err := r.Conn.Exec(query, event.EventID, event.EventType, event.UserID, event.EventTime, event.Payload)
	return err
}
