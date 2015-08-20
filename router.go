package neutron

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	r := Router{mux.NewRouter()}
	return &r
}

func (r *Router) Run(addr string) error {
	http.Handle("/", r)
	return http.ListenAndServe(addr, nil)
}

func (r *Router) AddFunc(path string, method string, f func(*Context)) {
	r.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			context := Context{w, r}
			f(&context)
		}).Methods(method)
}

func (r *Router) Get(path string, f func(*Context)) {
	r.AddFunc(path, "GET", f)
}

func (r *Router) Post(path string, f func(*Context)) {
	r.AddFunc(path, "POST", f)
}

