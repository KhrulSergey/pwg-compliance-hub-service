package config

import "github.com/caarlos0/env"

// AppConfig contains needed envs to run service
type AppConfig struct {
	AppVersion         string `env:"APP_VERSION" envDefault:"0.1-beta"`
	ServiceExternalUrl string `env:"COMPLIANCE_SERVICE_EXTERNAL_URL" envDefault:"localhost:7005"`
	ServiceProtocol    string `env:"COMPLIANCE_SERVICE_PROTOCOL" envDefault:"http"`
	Host               string `env:"COMPLIANCE_SERVICE_HOST" envDefault:"localhost"`
	Port               string `env:"COMPLIANCE_SERVICE_PORT" envDefault:"7005"`
	GrpcHost           string `env:"COMPLIANCE_SERVICE_HOST" envDefault:"localhost"`
	GrpcPort           string `env:"COMPLIANCE_SERVICE_GRPC_PORT" envDefault:"7105"`
	TokenSecret        string `env:"TOKEN_SECRET" envDefault:"24ed47b3b99144c2817fc788bad7b003"`
	LoggerMode         string `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

// InitAppConfig parses envs and constructs the config
func InitAppConfig() (*AppConfig, error) {
	var config AppConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
