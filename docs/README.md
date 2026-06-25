# dns-kong — docs

**Kong Gateway.** API-gateway services/routes/plugins via the Kong Admin API (Supabase-friendly).

## Install

```bash
togo install togo-framework/dns-kong
```

Registers on the [`dns`](https://github.com/togo-framework/dns) base; select it with **dns.provider in togo.yaml (or DNS_DRIVER)**, then use **`togo proxy`**.

## Interface

`Provider` — `UpsertRecord`/`DeleteRecord`/`ListRecords`, `UpsertProxyHost`/`DeleteProxyHost`, `UpsertRoute`/`DeleteRoute`.

## Configuration

| Env var | Description |
|---|---|
| `KONG_ADMIN` | Kong Admin API URL, e.g. `http://localhost:8001` (required). |
| `KONG_ADMIN_TOKEN` | Kong Admin API token, if RBAC is enabled. Optional. |

## Usage & notes

Creates service+route pairs and per-route plugins. Point a route at a Supabase backend and attach auth plugins to front Supabase with Kong.

## Example

```bash
togo proxy:host:add app.example.com http://localhost:3000 --provider kong --dry-run
```

## Links

- [Kong Admin API](https://docs.konghq.com/gateway/latest/admin-api/)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/dns-kong)
