package model

import (
    "context"
    "errors"

    "github.com/llaoj/oauth2nsso/config"
    "github.com/llaoj/oauth2nsso/pkg/ldap"
)

type User struct {
    ID       int    `gorm:"primary_key" json:"id"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

func (u *User) TableName() string {
    return "user"
}

func (u *User) Authentication(ctx context.Context, username, password string) (userID string, err error) {

    if config.Get().AuthMode == "ldap" {
        userID, err = ldap.UserAuthentication(username, password)
        return
    }

    if config.Get().AuthMode == "db" {
        // write your own user authentication logic
        // like:
        //   DB().WithContext(ctx).Where("name = ? AND password = ?", username, password).First(u)
        //   userID = u.ID
        if username != "user01" || password != "password01" {
            // test account:
            //   user01 password01userID = "user01"
            err = errors.New("用户名密码错误")
            return
        }

        userID = username
        return
    }

    return
}
