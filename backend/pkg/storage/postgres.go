package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

const (
	_defaultConnectionAttempts = 10
	_defaultConnectionTimeout  = time.Second

	_maxConnections        = int32(800)
	_minConnections        = int32(50)
	_maxConnectionLifeTime = time.Second * 30
	_maxIdleLifeTime       = time.Second * 60
)

type Postgres interface {
	Stats() *pgxpool.Stat
	Query(query string, args ...any) (pgx.Rows, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(query string, args ...interface{}) pgx.Row
}

type Pool struct {
	db *pgxpool.Pool
}

func Connect(connectionUrl string) (Postgres, error) {
	connectionAttempts := _defaultConnectionAttempts
	var result *pgxpool.Pool
	var err error
	connectionUrl += fmt.Sprintf(" pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%v pool_max_conn_idle_time=%v",
		_maxConnections, _minConnections, _maxConnectionLifeTime, _maxIdleLifeTime)
	for connectionAttempts > 0 {
		result, err = pgxpool.Connect(context.Background(), connectionUrl)
		if err == nil {
			break
		}

		log.Printf("ATTEMPT %d TO CONNECT TO POSTGRES BY URL %s FAILED: %s\n", connectionAttempts, connectionUrl, err.Error())

		connectionAttempts--

		time.Sleep(_defaultConnectionTimeout)
	}

	if result == nil {
		log.Printf("POSTGRES CONNECTION(%s) ERROR: %s\n", connectionUrl, err.Error())
		return nil, err
	}

	return &Pool{db: result}, nil
}

func (p Pool) Stats() *pgxpool.Stat {
	return p.db.Stat()
}

func (p Pool) Query(query string, args ...any) (pgx.Rows, error) {
	return p.db.Query(context.Background(), query, args[:]...)
}

func (p Pool) Get(dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(context.Background(), query, args[:]...)
	if err != nil {
		return err
	}
	return pgxscan.DefaultAPI.ScanOne(dest, rows)
}

func (p Pool) Select(dest interface{}, query string, args ...interface{}) error {
	rows, err := p.db.Query(context.Background(), query, args[:]...)
	if err != nil {
		return err
	}
	return pgxscan.DefaultAPI.ScanAll(dest, rows)
}

func (p Pool) Exec(query string, args ...any) (pgconn.CommandTag, error) {
	return p.db.Exec(context.Background(), query, args[:]...)
}

func (p Pool) QueryRow(query string, args ...interface{}) pgx.Row {
	return p.db.QueryRow(context.Background(), query, args[:]...)
}
