// Package kong is a togo dns driver for the Kong API gateway, driven through the
// Kong Admin API. Gateway routes (and proxy hosts, modeled as routes) become a
// Kong service + route pair. Pairs well with Supabase: point routes at the
// Supabase Kong/Edge endpoints to expose them behind your own gateway. DNS
// records return ErrUnsupported.
//
// Install: `togo install togo-framework/dns-kong`, set DNS_DRIVER=kong.
package kong

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/togo-framework/dns"
	"github.com/togo-framework/togo"
)

func init() {
	dns.RegisterDriver("kong", func(k *togo.Kernel) (dns.Provider, error) {
		admin := os.Getenv("KONG_ADMIN")
		if admin == "" {
			admin = "http://localhost:8001"
		}
		return &provider{
			admin: strings.TrimRight(admin, "/"),
			token: os.Getenv("KONG_ADMIN_TOKEN"),
			hc:    &http.Client{Timeout: 15 * time.Second},
		}, nil
	})
}

type provider struct {
	admin, token string
	hc           *http.Client
}

func (p *provider) do(ctx context.Context, method, path string, body any, out any) error {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, _ := http.NewRequestWithContext(ctx, method, p.admin+path, rdr)
	req.Header.Set("Content-Type", "application/json")
	if p.token != "" {
		req.Header.Set("Kong-Admin-Token", p.token)
	}
	resp, err := p.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return fmt.Errorf("dns-kong: %s %s -> %d: %s", method, path, resp.StatusCode, string(raw))
	}
	if out != nil && len(raw) > 0 {
		return json.Unmarshal(raw, out)
	}
	return nil
}

// upsertService PUTs a named service (idempotent) pointing at the upstream URL.
func (p *provider) upsertService(ctx context.Context, name, upstream string) error {
	return p.do(ctx, http.MethodPut, "/services/"+url.PathEscape(name), map[string]any{"url": upstream}, nil)
}

func (p *provider) upsertRoute(ctx context.Context, svc, domain, path string, plugins map[string]any) (string, error) {
	name := "togo_" + strings.NewReplacer(".", "_", "/", "_").Replace(domain+path)
	payload := map[string]any{
		"name":    name,
		"hosts":   []string{domain},
		"service": map[string]any{"name": svc},
	}
	if path != "" && path != "/" {
		payload["paths"] = []string{path}
	}
	var out struct {
		ID string `json:"id"`
	}
	// PUT by name = upsert.
	if err := p.do(ctx, http.MethodPut, "/routes/"+url.PathEscape(name), payload, &out); err != nil {
		return "", err
	}
	for plugin, cfg := range plugins {
		_ = p.do(ctx, http.MethodPost, "/routes/"+url.PathEscape(name)+"/plugins",
			map[string]any{"name": plugin, "config": cfg}, nil)
	}
	if out.ID != "" {
		return out.ID, nil
	}
	return name, nil
}

func (p *provider) UpsertRoute(ctx context.Context, rt dns.Route) (string, error) {
	svc := "togo_" + strings.NewReplacer(".", "_", "/", "_").Replace(rt.Domain+rt.Path)
	if err := p.upsertService(ctx, svc, rt.Upstream); err != nil {
		return "", err
	}
	return p.upsertRoute(ctx, svc, rt.Domain, rt.Path, rt.Plugins)
}

func (p *provider) UpsertProxyHost(ctx context.Context, h dns.ProxyHost) (string, error) {
	return p.UpsertRoute(ctx, dns.Route{Domain: h.Domain, Path: "/", Upstream: h.Upstream})
}

func (p *provider) DeleteRoute(ctx context.Context, id string) error {
	return p.do(ctx, http.MethodDelete, "/routes/"+url.PathEscape(id), nil, nil)
}
func (p *provider) DeleteProxyHost(ctx context.Context, id string) error {
	return p.DeleteRoute(ctx, id)
}

func (p *provider) UpsertRecord(context.Context, string, dns.Record) (string, error) {
	return "", dns.ErrUnsupported
}
func (p *provider) DeleteRecord(context.Context, string, string) error { return dns.ErrUnsupported }
func (p *provider) ListRecords(context.Context, string) ([]dns.Record, error) {
	return nil, dns.ErrUnsupported
}
