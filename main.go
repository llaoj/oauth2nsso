package main

import (
    "encoding/json"
    "net/http"
    "net/url"
    "html/template"
    "time"
    // "fmt"

    "github.com/dgrijalva/jwt-go"
    "github.com/go-session/session"
    "gopkg.in/oauth2.v3/errors"
    "gopkg.in/oauth2.v3/generates"
    "gopkg.in/oauth2.v3/manage"
    "gopkg.in/oauth2.v3/models"
    "gopkg.in/oauth2.v3/server"
    "gopkg.in/oauth2.v3/store"

    "oauth2/utils/yaml"
    "oauth2/utils/log"
    "oauth2/model"

)

type LoginPageDate struct {
    Client yaml.Client
    Scope [] yaml.Scope
    Error string
}

var srv *server.Server

func main() {
    yaml.Setup()
    log.Setup()
    model.Setup()

    //manager config
    manager := manage.NewDefaultManager()
    manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
    //token store
    manager.MustTokenStorage(store.NewMemoryTokenStore())
    //generate jwt access token
    manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
    clientStore := store.NewClientStore()
    for _, v := range yaml.Cfg.OAuth2.Client {
        clientStore.Set(v.ID, &models.Client{
            ID:     v.ID,
            Secret: v.Secret,
            Domain: v.Domain,
        })
    }
    manager.MapClientStorage(clientStore)
    // config oauth2 server
    srv = server.NewServer(server.NewConfig(), manager)
    srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
    srv.SetUserAuthorizationHandler(userAuthorizeHandler)
    srv.SetInternalErrorHandler(internalErrorHandler)
    srv.SetResponseErrorHandler(responseErrorHandler)

    // http server
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/authorize", authorizeHandler)
    http.HandleFunc("/token", tokenHandler)
    http.HandleFunc("/test", testHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    log.App.Info("Server is running at 9096 port.")
    err := http.ListenAndServe(":9096", nil)
    if err != nil {
        log.App.Error(err.Error())
    }
}

func passwordAuthorizationHandler(username, password string) (userID string, err error) {
    //自己实现验证逻辑
    var user model.User
    userID = user.GetUserIDByPwd(username, password)

    return
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
    store, err := session.Start(nil, w, r)
    if err != nil {
        return
    }

    uid, ok := store.Get("LoggedInUserID")
    if !ok {
        if r.Form == nil {
            r.ParseForm()
        }
        store.Set("RequestForm", r.Form)
        store.Save()

        w.Header().Set("Location", "/login")
        w.WriteHeader(http.StatusFound)
        return
    }

    userID = uid.(string)
    // 不记住用户
    // store.Delete("LoggedInUserID")
    // store.Save()

    return
}

func internalErrorHandler(err error) (re *errors.Response) {
    log.App.Error("Internal Error:", err.Error())
    return
}

func responseErrorHandler(re *errors.Response) {
    log.App.Error("Response Error:", re.Error.Error())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    store, err := session.Start(nil, w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    var pageData LoginPageDate
    if v, ok := store.Get("RequestForm"); ok {
        form := v.(url.Values)
        clientID := form.Get("client_id")
        for _, v := range yaml.Cfg.OAuth2.Client {
            if v.ID == clientID {
                pageData.Client = v
            }
        }
    }

    if r.Method == "POST" {
        if r.Form == nil {
            if err := r.ParseForm(); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
        userID := ""

        //账号密码验证
        if r.Form.Get("type") == "password" {
            //自己实现验证逻辑
            var user model.User
            userID = user.GetUserIDByPwd(r.Form.Get("username"), r.Form.Get("password"))
            if userID == "" {
                t, err := template.ParseFiles("tpl/login.html")
                if err != nil{
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
                pageData.Error = "用户名密码错误!"
                t.Execute(w, pageData)
            }
        }

        //扫码验证
        //手机验证码验证

        store.Set("LoggedInUserID", userID)
        store.Save()
        w.Header().Set("Location", "/authorize")
        w.WriteHeader(http.StatusFound)
        
        return
    }

    t, err := template.ParseFiles("tpl/login.html")
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, pageData)
}

// 首先进入执行
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
    store, err := session.Start(nil, w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    var form url.Values
    if v, ok := store.Get("RequestForm"); ok {
        form = v.(url.Values)
    }
    r.Form = form

    store.Delete("RequestForm")
    store.Save()

    err = srv.HandleAuthorizeRequest(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
    err := srv.HandleTokenRequest(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func testHandler(w http.ResponseWriter, r *http.Request) {
    token, err := srv.ValidationBearerToken(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data := map[string]interface{}{
        "expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
        "client_id": token.GetClientID(),
        "user_id": token.GetUserID(),
        "scope": token.GetScope(),
    }
    e := json.NewEncoder(w)
    e.SetIndent("", "  ")
    e.Encode(data)
}