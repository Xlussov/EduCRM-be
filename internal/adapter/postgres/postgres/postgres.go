package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DBName   string `yaml:"db_name" env:"POSTGRES_DB_NAME" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type Pool struct {
	conn *pgxpool.Pool
}

func New(ctx context.Context, cfg Config, log Logger) (*Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Pool{conn: pool}, nil
}

func (p *Pool) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Pool) Conn() *pgxpool.Pool {
	return p.conn
}
