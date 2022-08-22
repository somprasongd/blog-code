package main

import (
	"goapi-handlefunc/context"
	"goapi-handlefunc/handler"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	BASE_URL = "/api/%s"
	PORT     = ":8080"
)

func main() {
	r := mux.NewRouter()
	// use gorilla/mux
	setRouter(r)

	http.ListenAndServe(PORT, r)
}

func setRouter(r *mux.Router) {
	h := handler.TodoHandler{}

	todos := r.PathPrefix(BASE_URL + "/todos").Subrouter()
	todos.HandleFunc("", context.WrapHTTPHandler(h.CreateHandler)).Methods("POST")
	todos.HandleFunc("", context.WrapHTTPHandler(h.ListHandler)).Methods("GET")
	todos.HandleFunc("/{id:[0-9]+}", context.WrapHTTPHandler(h.GetHandler)).Methods("GET")
	todos.HandleFunc("/{id:[0-9]+}", context.WrapHTTPHandler(h.StatusUpdateHandler)).Methods("PATCH")
	todos.HandleFunc("/{id:[0-9]+}", context.WrapHTTPHandler(h.DeleteHandler)).Methods("DELETE")
}
