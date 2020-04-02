package model

import(
    "strconv"
    "crypto/sha256"
    "crypto/md5"
    "encoding/hex"
)

type User struct {
    UID int `gorm:"primary_key" json:"id"`
    UserName string `json:"username"`
    Pass string `json:"pass"`
    Tel string `json:"tel"`
    Email string `json:"email"`
    Salt string `json:"salt"`
}

func (u *User) TableName() string {
    return "admin_user"
}

func (u *User) GetUserIDByPwd(email, pwd string) (userID string) {
    userID = ""
    db.Where("tel=?", email).Or("email=?", email).First(u)
    if u.UID > 0 {
        bytes := sha256.Sum256([]byte(pwd))
        hash := hex.EncodeToString(bytes[:])
        hash = hash + u.Salt
        bytes2 := md5.Sum([]byte(hash))
        hash = hex.EncodeToString(bytes2[:])
        if hash == u.Pass {
            userID = strconv.Itoa(u.UID)
        }
    }

    return
}