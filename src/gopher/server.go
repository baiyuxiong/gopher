/*
读取配置文件,设置URL,启动服务器
*/

package gopher

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

func handlerFun(route Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("panic:", e)
				}
			}()*/

		handler := NewHandler(w, r)
		defer handler.Session.Close()

		url := r.Method + " " + r.URL.Path
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), url)
		if route.Permission == Everyone {
			route.HandlerFunc(handler)
		} else if route.Permission == Authenticated {
			_, ok := currentUser(handler)

			if !ok {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}

			route.HandlerFunc(handler)
		} else if route.Permission == Administrator {
			user, ok := currentUser(handler)
			if !ok {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			if !user.IsSuperuser {
				message(handler, "没有权限", "对不起，你没有权限进行该操作", "error")
				return
			}
			route.HandlerFunc(handler)
		}
	}
}

func StartServer() {
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.Handle("/get/package", websocket.Handler(getPackageHandler))
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	//http.Handle("/auth/signup", githubHandler)
	r := mux.NewRouter()
	for _, route := range routes {
		r.HandleFunc(route.URL, handlerFun(route))
	}

	http.Handle("/", r)

	fmt.Println("Server start on:", Config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil))
}
