package foundation

import (
	"context"
	"fmt"
	"strings"

	"github.com/eangle9/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb" // go-migrate needs it
	_ "github.com/golang-migrate/migrate/v4/source/file"          // go-migrate needs it
	"go.uber.org/zap"
)

func InitiateMigration(path, conn string, log log.Logger) *migrate.Migrate {
	url := fmt.Sprintf("cockroach://%v", strings.Split(conn, "://")[1])
	m, err := migrate.New(fmt.Sprintf("file://%v", path), url)
	if err != nil {
		log.Fatal(context.Background(), "failed to create migrator", zap.Error(err))
	}
	return m
}

func UpMigration(m *migrate.Migrate, log log.Logger) {
	err := m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(context.Background(), "could not migrate", zap.Error(err))
	}
}
