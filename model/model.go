package model

import (
    "fmt"
    "log"
    "time"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "github.com/llaoj/oauth2/config"
)

func DB() (db *gorm.DB) {
    var err error
    cfg := config.Get()
    db, err = gorm.Open(cfg.DB.Default.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
        cfg.DB.Default.User,
        cfg.DB.Default.Password,
        cfg.DB.Default.Host,
        cfg.DB.Default.DBName))
    if err != nil {
        log.Println(err)
    }

    db.SingularTable(true)
    db.LogMode(true)
    db.DB().SetMaxIdleConns(10)
    db.DB().SetMaxOpenConns(100)

    return
}

// func CloseDB() {
//     defer db.Close()
// }

type Model struct {
    ID        uint       `gorm:"primary_key" json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at"`
}
