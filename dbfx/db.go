package dbfx

import (
	"fmt"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Dsn string `mapstructure:"dsn" validate:"required,uri"`
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
