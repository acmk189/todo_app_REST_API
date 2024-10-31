package main

import (
	"net/http"

	"github.com/acmk189/todo_app_REST_API/handler"
	"github.com/acmk189/todo_app_REST_API/store"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析エラーを回避するため戻り値を無視
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	v := validator.New()
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/todos", lt.ServeHTTP)
	return mux
}
