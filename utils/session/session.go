package session

import (
    // "log"
    "net/http"
    "net/url"

    "encoding/gob"
    // "github.com/gorilla/sessions"
	"gopkg.in/boj/redistore.v1"

	"oauth2/utils/yaml"
)

// var store *sessions.CookieStore
var store *redistore.RediStore 

func Setup(){
    gob.Register(url.Values{})

    // store = sessions.NewCookieStore([]byte("SESSION_KEY"))
    store, _ = redistore.NewRediStore(yaml.Cfg.Redis.Default.Db, "tcp", yaml.Cfg.Redis.Default.Addr, "", []byte("secret-key"))
    // if err != nil {
    //     log.Fatal(err)

    //     return
    // }
}

func Get(r *http.Request, name string) (val interface{}, err error) {
    // Get a session.
    session, err := store.Get(r, yaml.Cfg.Session.Name)
    if err != nil {
        return
    }
    
    val = session.Values[name]

    return
}

func Set(w http.ResponseWriter, r *http.Request, name string, val interface{}) (err error) {
    // Get a session.
    session, err := store.Get(r, yaml.Cfg.Session.Name)
    if err != nil {
        return
    }

    session.Values[name] = val
    err = session.Save(r, w)

    return
}

func Delete(w http.ResponseWriter, r *http.Request, name string) (err error) {
    // Get a session.
    session, err := store.Get(r, yaml.Cfg.Session.Name)
    if err != nil {
        return
    }

    delete(session.Values, name)
    err = session.Save(r, w)

    return
}
