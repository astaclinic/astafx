package dbfx

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	UserName string `mapstructure:"user_name" yaml:"user_name" validate:"required"`
	Password string `mapstructure:"password" yaml:"password" validate:"required"`
	Host     string `mapstructure:"host" yaml:"host" validate:"required"`
	Database string `mapstructure:"database" yaml:"database" validate:"required"`
	Port     string `mapstructure:"port" yaml:"port" validate:"required"`
}

func init() {
	// config must have a default value for viper to load config from env variables
	// default value of empty string (zero value) will not pass the "required" config validation
	viper.SetDefault("postgres.dsn", "")
}

type Params struct {
	fx.In
	Config     *PostgresConfig
	GormLogger *GormLogger
}

func New(p Params) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", p.Config.UserName, p.Config.Password, p.Config.Host, p.Config.Port, p.Config.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: p.GormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("fail to initialize database: %w", err)
	}
	return db, nil
}

var Module = fx.Options(
	fx.Provide(New),
	fx.Provide(NewGormLogger),
	fx.Invoke(SetupGormPrometheus),
)
