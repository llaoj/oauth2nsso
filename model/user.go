package model

import "context"

type User struct {
    ID       int    `gorm:"primary_key" json:"id"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

func (u *User) TableName() string {
    return "user"
}

func (u *User) GetUserIDByPwd(ctx context.Context, username, password string) (userID string) {
    // use the db conn
    // write your own user authentication logic
    // like:
    // db.WithContext(ctx).Where("name = ? AND password = ?", username, password).First(u)
    // userID = u.ID

    // test account: admin admin
    if username == "admin" && password == "admin" {
        userID = "admin"
    }

    return
}
