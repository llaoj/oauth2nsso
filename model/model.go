package model

import (
    "fmt"
    "log"
    "time"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "github.com/llaoj/oauth2/config"
)

var db *gorm.DB

func Setup() {
    var err error
    cfg := config.Get()
    db, err = gorm.Open(cfg.Db.Default.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
        cfg.Db.Default.User,
        cfg.Db.Default.Password,
        cfg.Db.Default.Host,
        cfg.Db.Default.DbName))
    if err != nil {
        log.Println(err)
    }

    db.SingularTable(true)
    db.LogMode(true)
    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
    defer db.Close()
}

type Model struct {
    ID        uint       `gorm:"primary_key" json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at"`
}
