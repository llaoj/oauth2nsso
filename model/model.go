package model

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"

    "github.com/llaoj/oauth2nsso/config"
)

var db *gorm.DB

func DB() *gorm.DB {
    if db != nil {
        return db
    }

    var err error
    cfg := config.Get().DB.Default

    switch cfg.Type {
    case "mysql":
        // dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
        dsn := fmt.Sprintf(
            "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
            cfg.User,
            cfg.Password,
            cfg.Host,
            cfg.Port,
            cfg.DBName)
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Silent),
        })
    case "postgresql":
        // to do something
        // ...
    }

    if err != nil {
        log.Fatal(err)
    }

    // GORM 使用 database/sql 维护连接池
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal(err)
    }
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    return db
}

type Model struct {
    ID        uint       `gorm:"primary_key" json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at"`
}
