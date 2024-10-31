package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acmk189/todo_app_REST_API/entity"
	"github.com/acmk189/todo_app_REST_API/store"
	"github.com/acmk189/todo_app_REST_API/testutil"
	"github.com/go-playground/validator/v10"
)

// TestAddTask は、さまざまなリクエストシナリオに対する AddTask ハンドラーをテストします。
// テーブル駆動テストを使用して、異なるリクエストケースに対するハンドラーのレスポンスステータスとボディを
// 期待される値と比較して検証します。
//
// テストケースには以下が含まれます:
// - "ok": HTTP 200 OK ステータスと期待されるレスポンスボディを返すべき有効なリクエスト。
// - "bad_request": HTTP 400 Bad Request ステータスと期待されるレスポンスボディを返すべき無効なリクエスト。
//
// 各テストケースはリクエストと期待されるレスポンスデータをゴールデンファイルから読み込み、
// AddTask ハンドラーにリクエストを送信し、レスポンスのステータスとボディを検証します。
func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"bad_request": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			sut := &AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			// レスポンスのステータスコードとボディを検証
			testutil.AssertResponce(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
