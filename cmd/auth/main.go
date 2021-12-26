package main

import (
	"fullstack-shop/pkg/auth/model"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/reg", model.Register)
	mux.HandleFunc("/login", model.Login)
	mux.Handle("/getme", model.Middleware(http.HandlerFunc(model.GetMe)))
	mux.Handle("/logout", model.Middleware(http.HandlerFunc(model.Logout)))
	mux.HandleFunc("/refresh", model.Refresh)
	http.ListenAndServe(":8080", mux)
}
