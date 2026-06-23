package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HMasataka/beyond/internal/handler"
	"github.com/HMasataka/beyond/internal/openapi"
)

type stubPinger struct{}

func (stubPinger) PingContext(_ context.Context) error { return nil }

func TestServerHealthz(t *testing.T) {
	// 集約した Server が生成インターフェースを満たし、ルータ経由で応答することを検証する。
	ts := httptest.NewServer(openapi.Handler(openapi.NewStrictHandler(handler.NewServer(stubPinger{}), nil)))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	var got openapi.HealthStatus
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got.Status != openapi.HealthStatusStatusOk {
		t.Errorf("status = %q, want %q", got.Status, openapi.HealthStatusStatusOk)
	}
}
