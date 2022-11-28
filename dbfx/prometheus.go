package dbfx

import (
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

func SetupGormPrometheus(db *gorm.DB) {
	db.Use(prometheus.New(prometheus.Config{
		StartServer: false,
	}))
}
