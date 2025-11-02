package infrastructure

import (
	"fmt"

	"backend/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// NewMySQLConnection establishes a gorm DB connection using environment configuration.
func NewMySQLConnection() (*gorm.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	dbCfg := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	return gorm.Open("mysql", dsn)
}
