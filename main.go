package main

import (
    // "fmt"
    "log"
    "time"
    "net/url"
    "net/http"
    "encoding/json"
    "html/template"

    "github.com/dgrijalva/jwt-go"
    "gopkg.in/oauth2.v3/errors"
    "gopkg.in/oauth2.v3/generates"
    "gopkg.in/oauth2.v3/manage"
    "gopkg.in/oauth2.v3/models"
    "gopkg.in/oauth2.v3/server"
    "gopkg.in/oauth2.v3/store"
    "github.com/go-redis/redis"
    oredis "gopkg.in/go-oauth2/redis.v3"

    "oauth2/model"
    "oauth2/config"
    "oauth2/pkg/session"
)

var srv *server.Server
var mgr *manage.Manager

func main() {
    time.Sleep(30 * time.Second)
    config.Setup()

    // init db connection
    // configure db in app.yaml then uncomment
    // model.Setup()

    session.Setup()

    // manager config
    mgr = manage.NewDefaultManager()
    mgr.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
    // token store
    // manager.MustTokenStorage(store.NewMemoryTokenStore())
    // use redis token store
    mgr.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
        Addr: config.Get().Redis.Default.Addr,
        DB: config.Get().Redis.Default.Db,
    }))

    // access token generate method: jwt
    mgr.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
    clientStore := store.NewClientStore()
    for _, v := range config.Get().OAuth2.Client {
        clientStore.Set(v.ID, &models.Client{
            ID:     v.ID,
            Secret: v.Secret,
            Domain: v.Domain,
        })
    }
    mgr.MapClientStorage(clientStore)
    // config oauth2 server
    srv = server.NewServer(server.NewConfig(), mgr)
    srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
    srv.SetUserAuthorizationHandler(userAuthorizeHandler)
    srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
    srv.SetInternalErrorHandler(internalErrorHandler)
    srv.SetResponseErrorHandler(responseErrorHandler)

    // http server
    http.HandleFunc("/authorize", authorizeHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/logout", logoutHandler)
    http.HandleFunc("/token", tokenHandler)
    http.HandleFunc("/test", testHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    log.Println("Server is running at 9096 port.")
    log.Fatal(http.ListenAndServe(":9096", nil))
}

func passwordAuthorizationHandler(username, password string) (userID string, err error) {
    var user model.User
    userID = user.GetUserIDByPwd(username, password)

    return
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
    v, _ := session.Get(r, "LoggedInUserID")
    if v == nil {
       if r.Form == nil {
            r.ParseForm()
        }
        session.Set(w, r, "RequestForm", r.Form)
        
        w.Header().Set("Location", "/login")
        w.WriteHeader(http.StatusFound)

        return
    }
    userID = v.(string)

    // 不记住用户
    // store.Delete("LoggedInUserID")
    // store.Save()

    return
}

// 根据client注册的scope
// 过滤非法scope
func authorizeScopeHandler (w http.ResponseWriter, r *http.Request) (scope string, err error) {
   if r.Form == nil {
        r.ParseForm()
    }
    s := config.ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
    if s == nil {
        http.Error(w, "Invalid Scope", http.StatusBadRequest)
        return
    }
    scope = config.ScopeJoin(s)

    return
}

func internalErrorHandler(err error) (re *errors.Response) {
    log.Println("Internal Error:", err.Error())
    return
}

func responseErrorHandler(re *errors.Response) {
    log.Println("Response Error:", re.Error.Error())
}

// 首先进入执行
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
    var form url.Values
    if v, _ := session.Get(r, "RequestForm"); v != nil {
       r.ParseForm()
        if r.Form.Get("client_id") == "" {
            form = v.(url.Values)
        } 
    }
    r.Form = form

    if err := session.Delete(w, r, "RequestForm"); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := srv.HandleAuthorizeRequest(w, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }
}

type TplData struct {
    Client config.Client
    // 用户申请的合规scope
    Scope []config.Scope
    Error string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    form, err := session.Get(r, "RequestForm")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if form == nil {
        http.Error(w, "Invalid Request", http.StatusBadRequest)
        return
    }
    clientID := form.(url.Values).Get("client_id")
    scope := form.(url.Values).Get("scope")

    // 页面数据
    data := TplData{
        Client: config.GetClient(clientID),
        Scope: config.ScopeFilter(clientID, scope),
    }
    if data.Scope == nil {
        http.Error(w, "Invalid Scope", http.StatusBadRequest)
        return
    }

    if r.Method == "POST" {
        if r.Form == nil {
            if err := r.ParseForm(); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }
        var userID string

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
                data.Error = "用户名密码错误!"
                t.Execute(w, data)

                return
            }
        }

        //扫码验证
        //手机验证码验证

        if err := session.Set(w, r, "LoggedInUserID", userID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Location", "/authorize")
        w.WriteHeader(http.StatusFound)
        
        return
    }

    t, err := template.ParseFiles("tpl/login.html")
    if err != nil{
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    t.Execute(w, data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    if r.Form == nil {
        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    redirectURI := r.Form.Get("redirect_uri")
    if _, err := url.Parse(redirectURI); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    if err := session.Delete(w, r, "LoggedInUserID"); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Location", redirectURI)
    w.WriteHeader(http.StatusFound)
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
    cli, err := mgr.GetClient(token.GetClientID())
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data := map[string]interface{}{
        "expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
        "user_id": token.GetUserID(),
        "client_id": token.GetClientID(),
        "scope": token.GetScope(),
        "domain": cli.GetDomain(),
    }
    e := json.NewEncoder(w)
    e.SetIndent("", "  ")
    e.Encode(data)
}