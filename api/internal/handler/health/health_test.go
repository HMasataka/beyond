package health_test

import (
	"context"
	"testing"

	"github.com/HMasataka/beyond/internal/handler/health"
	"github.com/HMasataka/beyond/internal/openapi"
)

func TestGetHealthz(t *testing.T) {
	got, err := health.New().GetHealthz(context.Background(), openapi.GetHealthzRequestObject{})
	if err != nil {
		t.Fatalf("GetHealthz returned error: %v", err)
	}
	res, ok := got.(openapi.GetHealthz200JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want GetHealthz200JSONResponse", got)
	}
	if res.Status != openapi.HealthStatusStatusOk {
		t.Errorf("status = %q, want %q", res.Status, openapi.HealthStatusStatusOk)
	}
}
