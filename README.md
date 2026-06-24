<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/dns-kong</h1>
  <p>Kong API-gateway driver (+ Supabase) for togo dns.</p>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/dns-kong"><img src="https://pkg.go.dev/badge/github.com/togo-framework/dns-kong.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/dns-kong
```
<!-- /togo-header -->

[Kong](https://konghq.com/) API-gateway driver for togo's
[`dns`](https://github.com/togo-framework/dns) subsystem, driven through the Kong
Admin API. Each route becomes a Kong **service + route** pair; gateway plugins
(rate-limiting, auth, …) are applied from `Route.Plugins`.

### Supabase

Kong is the gateway Supabase ships in front of its stack. Point routes at the
Supabase Kong/Edge endpoints to expose Supabase services behind your own gateway,
or run a standalone Kong in front of `supabase` plugins.

## Config

| Env | Meaning |
|-----|---------|
| `DNS_DRIVER` | set to `kong` |
| `KONG_ADMIN` | Admin API URL (default `http://localhost:8001`) |
| `KONG_ADMIN_TOKEN` | sent as `Kong-Admin-Token` (Kong Enterprise / RBAC) |

```go
svc, _ := dns.FromKernel(k)
svc.UpsertRoute(ctx, dns.Route{
    Domain:   "api.example.com",
    Path:     "/v1",
    Upstream: "http://my-service:9000",
    Plugins:  map[string]any{"rate-limiting": map[string]any{"minute": 60}},
})
```

DNS records return `dns.ErrUnsupported`.

<!-- togo-sponsors -->
---
<div align="center">
  <h3>Premium sponsors</h3>
  <p><a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp; <a href="https://one-studio.co"><strong>One Studio</strong></a></p>
  <p><sub>Support togo — <a href="https://github.com/sponsors/fadymondy">become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
