package ldap

import (
    "errors"
    "fmt"

    ldap "github.com/go-ldap/ldap/v3"
    "github.com/llaoj/oauth2/config"
)

func UserAuthentication(username, password string) (userID string, err error) {

    cfg := config.Get().LDAP

    l, err := ldap.DialURL(cfg.URL)
    if err != nil {
        return
    }
    defer l.Close()

    // Reconnect with TLS
    // err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
    // if err != nil {
    //     return
    // }

    // First bind with a read only user
    err = l.Bind(cfg.SearchDN, cfg.SearchPassword)
    if err != nil {
        return

    }

    // Search for the given username
    searchRequest := ldap.NewSearchRequest(
        cfg.BaseDN,
        ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
        fmt.Sprintf(cfg.Filter, ldap.EscapeFilter(username)),
        []string{"dn"},
        nil,
    )

    sr, err := l.Search(searchRequest)
    if err != nil {
        return
    }

    if len(sr.Entries) != 1 {
        err = errors.New("用户不存在或者不唯一")
        return
    }

    userdn := sr.Entries[0].DN

    // Bind as the user to verify their password
    err = l.Bind(userdn, password)
    if err != nil {
        return
    }

    // Rebind as the read only user for any further queries
    err = l.Bind(cfg.SearchDN, cfg.SearchPassword)
    if err != nil {
        return
    }

    userID = username
    return
}
