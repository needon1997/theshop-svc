package model

import (
	"database/sql"
	"fmt"
	"github.com/needon1997/theshop-svc/internal/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var DB *gorm.DB

func InitConnection() (err error) {
	conf := config.ServerConfig.MySqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.RestURI)
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return err
	}
	var sqlDB *sql.DB
	sqlDB, err = DB.DB()
	if err != nil {
		return err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(3)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	err = DB.AutoMigrate(&Category{}, &Brand{}, &Goods{}, &GoodsBrandCategory{}, &Banner{})
	return err
}
