package config

import (
        // "fmt"
        "strings"
)

func Get() *App {
    return &cfg
}

func GetClient(clientID string) (cli Client) {
    for _, v := range cfg.OAuth2.Client {
        if v.ID == clientID {
            cli = v
        }
    }

    return
}

func ScopeJoin(scope []Scope) string {
    var s []string
    for _, sc := range scope {
        s = append(s, sc.ID)
    }
    return strings.Join(s,",")
}

func ScopeFilter(clientID string, scope string) (s []Scope) {
    cli := GetClient(clientID)
    sl := strings.Split(scope, ",")
    for _, str := range sl {
        for _, sc := range cli.Scope {
            if str == sc.ID {
                s = append(s, sc)
            }
        }
    }

    return
}