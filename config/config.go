package config

import (
    // "fmt"
    // "strings"
    "flag"
    "io/ioutil"
    "log"

    "gopkg.in/yaml.v2"
)

var cfg App

func Setup() {
    path := flag.String("config", "/etc/oauth2/app.yaml", "config file's（app.yaml）absolute path")
    flag.Parse()

    content, err := ioutil.ReadFile(*path)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    err = yaml.Unmarshal(content, &cfg)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
}
