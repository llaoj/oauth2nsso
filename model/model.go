package model

import (
    "log"
    "fmt"
    "time"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"

    "oauth2/utils/yaml"
)

var db *gorm.DB

func Setup() {
    var err error
    db, err = gorm.Open(yaml.Cfg.Db.Default.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", 
        yaml.Cfg.Db.Default.User, 
        yaml.Cfg.Db.Default.Password, 
        yaml.Cfg.Db.Default.Host, 
        yaml.Cfg.Db.Default.DbName))
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
    ID        uint `gorm:"primary_key" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at"`
}
