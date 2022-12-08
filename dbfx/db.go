package dbfx

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Dsn string `mapstructure:"dsn" yaml:"dsn" validate:"required,uri"`
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
	db, err := gorm.Open(postgres.Open(p.Config.Dsn), &gorm.Config{
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
