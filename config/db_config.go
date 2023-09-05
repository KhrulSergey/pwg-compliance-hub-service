package config

import "github.com/caarlos0/env"

type DBConfig struct {
	DatabaseHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	DatabasePort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	DatabaseUser     string `env:"POSTGRES_USER" envDefault:"paywithglass_auth"`
	DatabasePassword string `env:"POSTGRES_PASSWORD" envDefault:"qwerty"`
	DatabaseName     string `env:"POSTGRES_DB" envDefault:"compliance_hub_service_database"`
	DatabaseRootCA   string `env:"POSTGRES_ROOT_CA" envDefault:""`
	ConnMaxLifeTime  int    `env:"POSTGRES_CONN_MAX_LIFETIME_S" envDefault:"3600"`
}

func NewDBConfig() (*DBConfig, error) {
	config := &DBConfig{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
