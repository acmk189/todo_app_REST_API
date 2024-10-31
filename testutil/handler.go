package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// JSONの比較を行うためのヘルパー関数
func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("failed to unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("failed to unmarshal want %q: %v", got, err)
	}
	if diff := cmp.Diff(jw, jg); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

// AssertResponceはHTTPレスポンスをアサートするテストヘルパー関数です。
// レスポンスのステータスコードが期待されるステータスと一致するか、レスポンスボディが期待されるボディと一致するかを確認します。
// レスポンスボディが空でない場合、AssertJSONを呼び出して期待されるJSONボディと実際のJSONボディを比較します。
//
// パラメータ:
//   - t: テストの失敗を報告するために使用されるテストオブジェクト。
//   - got: チェックされる実際のHTTPレスポンス。
//   - status: 期待されるHTTPステータスコード。
//   - body: 期待されるレスポンスボディ。
//
// この関数は、テストが終了した後にレスポンスボディが閉じられることも保証します。
func AssertResponce(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()

	t.Cleanup(func() { _ = got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if got.StatusCode != status {
		t.Errorf("want status %d, but got %d, body: %q", status, got.StatusCode, gb)
	}

	if len(gb) == 0 && len(body) == 0 {
		// 期待も実体もレスポンスボディが空なのでAssertJSONを呼び出さない
		return
	}
	AssertJSON(t, body, gb)
}

// 入力値と期待値をファイルから取得
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %q: %v", path, err)
	}
	return bt
}
