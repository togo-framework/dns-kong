package kong

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/togo-framework/dns"
)

func TestUpsertRouteCreatesServiceAndRoute(t *testing.T) {
	var svcPut, routePut bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/services/") {
			svcPut = true
		}
		if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/routes/") {
			routePut = true
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":"route-uuid"}`))
			return
		}
		w.WriteHeader(201)
	}))
	defer srv.Close()

	p := &provider{admin: srv.URL, hc: srv.Client()}
	id, err := p.UpsertRoute(context.Background(), dns.Route{Domain: "api.example.com", Path: "/v1", Upstream: "http://svc:9000"})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if !svcPut || !routePut {
		t.Fatalf("svcPut=%v routePut=%v", svcPut, routePut)
	}
	if id != "route-uuid" {
		t.Fatalf("id=%q", id)
	}
}
