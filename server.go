package main

import (
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"html/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var srv *server.Server

func main() {
	//manager config
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	//token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	//generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("00000000"), jwt.SigningMethodHS512))
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})
	manager.MapClientStorage(clientStore)

	// config oauth2 server
	srv = server.NewServer(server.NewConfig(), manager)
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)
	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetInternalErrorHandler(internalErrorHandler)
	srv.SetResponseErrorHandler(responseErrorHandler)

	// http server
	http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/consent", consentHandler)
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/token", tokenHandler)
	http.HandleFunc("/test", testHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func passwordAuthorizationHandler(username, password string) (userID string, err error) {
	if username == "test" && password == "test" {
		userID = "test"
	}
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
		store.Set("ReturnUri", r.Form)
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
	log.Println("Internal Error:", err.Error())
	return
}

func responseErrorHandler(re *errors.Response) {
	log.Println("Response Error:", re.Error.Error())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if r.Form.Get("password") != "test" {
			t, err := template.ParseFiles("tpl/login.html")
			if err != nil{
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		    }
    		data := struct{Error string}{"用户名密码错误 !"}
		    t.Execute(w,data)
		}

		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		// w.Header().Set("Location", "/consent")
		w.Header().Set("Location", "/authorize")
		w.WriteHeader(http.StatusFound)
		return
	}

	t, err := template.ParseFiles("tpl/login.html")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
    t.Execute(w,nil)
}

// func consentHandler(w http.ResponseWriter, r *http.Request) {
// 	store, err := session.Start(nil, w, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	if _, ok := store.Get("LoggedInUserID"); !ok {
// 		w.Header().Set("Location", "/login")
// 		w.WriteHeader(http.StatusFound)
// 		return
// 	}
// 	t, err := template.ParseFiles("tpl/consent.html")
// 	if err != nil{
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
//     }
//     t.Execute(w,nil)
// }

// 首先进入执行
func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	r.Form = form

	store.Delete("ReturnUri")
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
