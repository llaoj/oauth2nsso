package config

import (
        // "fmt"
        // "strings"
        "log"
        "io/ioutil"

        "gopkg.in/yaml.v2"
)

var cfg App

func Setup() {
    content, err := ioutil.ReadFile("app.yaml")
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    err = yaml.Unmarshal(content, &cfg)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
}