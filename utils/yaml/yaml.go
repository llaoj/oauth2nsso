package yaml

import (
        // "fmt"
        "log"
        "io/ioutil"

        "gopkg.in/yaml.v2"
)

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type App struct {
    Db struct {
        Default Db
    }
    Redis struct {
        Default Redis
    }
    OAuth2 struct {
        Client []Client `yaml:"client"`
    } `yaml:"oauth2"`
}

type Db struct {
    Type string `yaml:"type"`
    Host string `yaml:"host"`
    Port int `yaml:"port"`
    User string `yaml:"user"`
    Password string `yaml:"password"`
    DbName string `yaml:"dbname"`
}

type Redis struct {
    Addr string `yaml:"addr"`
    Password string `yaml:"password"`
    Db int `yaml:"db"`
}

type Client struct {
    ID string `yaml:"id"`
    Secret string `yaml:"secret"`
    Name string `yaml:"name"`
    Domain string `yaml:"domain"`
}

var Cfg App

func Setup() {
    content, err := ioutil.ReadFile("app.yaml")
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    err = yaml.Unmarshal(content, &Cfg)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
}