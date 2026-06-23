package readiness_test

import (
	"context"
	"errors"
	"testing"

	"github.com/HMasataka/beyond/internal/handler/readiness"
	"github.com/HMasataka/beyond/internal/openapi"
)

type fakePinger struct {
	err error
}

func (p fakePinger) PingContext(_ context.Context) error {
	return p.err
}

func TestGetReadyzReturnsOKWhenPingSucceeds(t *testing.T) {
	// Given: ping が成功する pinger
	h := readiness.New(fakePinger{err: nil})

	// When: GetReadyz を呼ぶ
	got, err := h.GetReadyz(context.Background(), openapi.GetReadyzRequestObject{})

	// Then: 200 応答で status は ok
	if err != nil {
		t.Fatalf("GetReadyz returned error: %v", err)
	}
	res, ok := got.(openapi.GetReadyz200JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want GetReadyz200JSONResponse", got)
	}
	if res.Status != openapi.ReadinessStatusStatusOk {
		t.Errorf("status = %q, want %q", res.Status, openapi.ReadinessStatusStatusOk)
	}
}

func TestGetReadyzReturnsUnavailableWhenPingFails(t *testing.T) {
	// Given: ping が失敗する pinger
	h := readiness.New(fakePinger{err: errors.New("connection refused")})

	// When: GetReadyz を呼ぶ
	got, err := h.GetReadyz(context.Background(), openapi.GetReadyzRequestObject{})

	// Then: 503 応答で status は unavailable
	if err != nil {
		t.Fatalf("GetReadyz returned error: %v", err)
	}
	res, ok := got.(openapi.GetReadyz503JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want GetReadyz503JSONResponse", got)
	}
	if res.Status != openapi.ReadinessStatusStatusUnavailable {
		t.Errorf("status = %q, want %q", res.Status, openapi.ReadinessStatusStatusUnavailable)
	}
}
