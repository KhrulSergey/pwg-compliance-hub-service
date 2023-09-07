package storage

import (
	"embed"
	"fmt"
	"gitlab.com/smdgroup/firmshift/backend/compliance-hub-service.git/config"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	//go:embed migrations/*.sql
	migrations embed.FS
)

func InitGormDB(config *config.DBConfig) (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s %s",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseName,
		sslConnectionParams(config.DatabaseRootCA),
	)
	conn, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: connectionString,
			},
		),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	db, err := conn.DB()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Second * time.Duration(config.ConnMaxLifeTime))
	return conn, err
}

func sslConnectionParams(dbRootCA string) string {
	if dbRootCA == "" {
		return "sslmode=disable"
	}
	return fmt.Sprintf("sslmode=verify-full sslrootcert=%s", dbRootCA)
}

func MigrateGorm(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:      "V202309041932540  create tables",
			Migrate: migrateFromFile("migrations/V202309041652540__create_tables.sql"),
		},
	})

	return m.Migrate()
}

func migrateFromFile(migrationFile string) gormigrate.MigrateFunc {
	return func(tx *gorm.DB) error {
		sql, err := migrations.ReadFile(migrationFile)
		if err != nil {
			return err
		}
		return tx.Exec(string(sql)).Error
	}
}
