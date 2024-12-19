package foundation

import (
	"context"
	"time"

	"github.com/eangle9/log"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitDB(url string, log log.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), "Failed to connect to the database", zap.Error(err))
	}
	idleConnTimeout := viper.GetDuration("database.idle_conn_timeout")
	if idleConnTimeout == 0 {
		idleConnTimeout = 4 * time.Minute
	}

	config.ConnConfig.Logger = log
	config.MaxConnIdleTime = idleConnTimeout
	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), "Failed to connect to database", zap.Error(err))
	}
	if _, err := conn.Exec(context.Background(), "show tables"); err != nil {
		log.Fatal(context.Background(), "Failed to ping database", zap.Error(err))
	}

	return conn
}
