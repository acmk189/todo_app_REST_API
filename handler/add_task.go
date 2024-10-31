package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/acmk189/todo_app_REST_API/entity"
	"github.com/acmk189/todo_app_REST_API/store"
	"github.com/go-playground/validator/v10"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

// http.HandlerFunc型を満たすメソッド
// 正常・異常時にRespondJSONでJSONを返す
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		// Validator.Structで検証する内容をタグに記述
		// 今回はTitleが必須なので以下を設定
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:     b.Title,
		Status:    entity.TaskStatusTodo,
		CreatedAt: time.Now(),
	}
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID int `json:"id"`
	}{ID: id}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
