package config

import (
	"fmt"
	"github.com/krez3f4l/audit_logger/internal/consts"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	env    string     `mapstructure:"env" validate:"required,oneof=local dev prod"`
	Server GRPCServer `mapstructure:"grpc_server" validate:"required"`
	DBConn DBConn     `mapstructure:"db_conn" validate:"required"`
}

type GRPCServer struct {
	Port int `mapstructure:"port" validate:"required,min=1,max=65535"`
}

type DBConn struct {
	URI      string `mapstructure:"uri" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Database string `mapstructure:"database" validate:"required"`
}

const (
	defaultConfigDir  = "configs"
	defaultConfigName = "config"
)

func NewConfig(configDir, configName string) (*Config, error) {
	v := viper.New()

	if configDir == "" {
		configDir = defaultConfigDir
	}
	if configName == "" {
		configName = defaultConfigName
	}

	v.AddConfigPath(configDir)
	v.SetConfigName(configName)
	v.SetConfigType(consts.CfgType)

	v.SetEnvPrefix(consts.EnvVarPrefix)
	v.AutomaticEnv()
	v.BindEnv("db_conn.username", consts.EnvVarPrefix+"_DB_CONN_USER")
	v.BindEnv("db_conn.password", consts.EnvVarPrefix+"_DB_CONN_PASSWORD")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		if strings.Contains(err.Error(), "DBConn.Username") {
			return nil, fmt.Errorf("database username must be set via %s_DB_CONN_USER env var", consts.EnvVarPrefix)
		}
		if strings.Contains(err.Error(), "DBConn.Password") {
			return nil, fmt.Errorf("database password must be set via %s_DB_CONN_PASSWORD env var", consts.EnvVarPrefix)
		}
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
