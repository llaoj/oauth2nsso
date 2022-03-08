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
    path := flag.String("config", "/etc/oauth2nsso/config.yaml", "the absolute path of config.yaml")
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
