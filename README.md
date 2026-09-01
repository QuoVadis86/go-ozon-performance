# Ozon Performance API Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](https://golang.org/dl/)
[![Go Reference](https://pkg.go.dev/badge/github.com/QuoVadis86/go-ozon-performance.svg)](https://pkg.go.dev/github.com/QuoVadis86/go-ozon-performance)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[📖 中文文档](README_CN.md)

A Go client library for the [Ozon Performance API](https://docs.ozon.ru/api/performance/).  
Covers **44 API methods** across **6 service modules** with **115 generated types**.

## Features

- Complete coverage of the Ozon Performance API v2 (OpenAPI 3.0)
- Automatic OAuth2 token management — fetches, caches and refreshes the Bearer token transparently
- Strongly-typed request/response structures generated from `swagger.json`
- `transport.Uint64` tolerant of both JSON number and string wire formats
- Async report flow support (submit → check → download CSV/ZIP)
- Enum types with documented values and descriptions
- Context support for cancellation and timeouts
- No external dependencies beyond the Go standard library
- Concurrent-safe transport with connection pooling

## Installation

```bash
go get github.com/QuoVadis86/go-ozon-performance
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    ozon "github.com/QuoVadis86/go-ozon-performance"
)

func main() {
    // Credentials from the seller cabinet: Settings → API keys → Performance API
    client := ozon.NewClient("your-client-id", "your-client-secret", nil)

    // List campaigns
    ctx := context.Background()
    campaigns, err := client.Campaign.ListCampaigns(ctx, &campaign.ListCampaignsRequest{
        AdvObjectType: "SKU",
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, c := range campaigns.List {
        fmt.Printf("id=%s title=%s state=%s\n", c.ID, c.Title, c.State)
    }
}
```

## API Coverage

| Service | Methods | Description |
|---------|---------|-------------|
| [Campaign](./campaign/service.go) | 5 | Campaign list, objects, bid limits, min bid by SKU |
| [Statistics](./statistics/service.go) | 17 | Campaign statistics reports, phrase stats, SKU stats |
| [Ad](./ad/service.go) | 5 | Campaign create, activate, deactivate, patch, dynamic budget |
| [Product](./product/service.go) | 5 | Product bids management in CPC campaigns |
| [Search-Promo](./searchpromo/service.go) | 8 | Pay-per-order campaigns, recommended bids |
| [Vendor](./vendor/service.go) | 4 | External traffic analysis reports |

## Examples

### List Campaigns

```go
import "github.com/QuoVadis86/go-ozon-performance/campaign"

resp, err := client.Campaign.ListCampaigns(ctx, &campaign.ListCampaignsRequest{
    AdvObjectType: "SKU",
    State:         campaign.StateCampaignStateRunning,
    Page:          1,
    PageSize:      100,
})
```

### Submit an Async Statistics Report

```go
import (
    "github.com/QuoVadis86/go-ozon-performance/statistics"
    "github.com/QuoVadis86/go-ozon-performance/transport"
)

// 1. Submit — returns a UUID
sub, err := client.Statistics.SubmitRequest(ctx, &statistics.StatisticsRequest{
    Campaigns: []transport.Uint64{48852},
    DateFrom:  "2026-08-01",
    DateTo:    "2026-08-31",
    GroupBy:   statistics.GroupByDate,
})

// 2. Poll until ready
for {
    st, err := client.Statistics.StatisticsCheck(ctx, sub.UUID)
    if err != nil {
        log.Fatal(err)
    }
    if st.State == statistics.StatisticsRequestStateOK {
        // 3. Download the CSV report
        data, err := client.Statistics.DownloadStatistics(ctx, &statistics.DownloadStatisticsRequest{
            UUID: sub.UUID,
        })
        // write data to file...
        break
    }
    time.Sleep(5 * time.Second)
}
```

### Manage Product Bids

```go
import (
    "github.com/QuoVadis86/go-ozon-performance/product"
    "github.com/QuoVadis86/go-ozon-performance/transport"
)

// Add products to a CPC campaign
err := client.Product.AddProducts(ctx, campaignID, &product.AddProductsRequest{
    Bids: []product.AddProductsRequestProduct{
        {
            SKU: transport.Uint64(48852),
            Bid: transport.Uint64(100000), // RUB in millionths
        },
    },
})
```

## Token Management

Token acquisition is automatic. `POST /api/client/token` is called once, cached across
concurrent callers, and refreshed before expiry. A `401` response triggers a forced
refresh with a single retry. For advanced use, the transport supports a `TokenStore`
to persist tokens across restarts:

```go
client := ozon.NewClient(clientID, secret, &ozon.ClientOptions{
    TokenStore: myPersistentStore, // implements transport.TokenStore
})
```

## Error Handling

All API error responses return `*transport.APIError` with a status code and message.
Check it with `errors.As`:

```go
import "github.com/QuoVadis86/go-ozon-performance/transport"

resp, err := client.Campaign.ListCampaigns(ctx, req)
if err != nil {
    var apiErr *transport.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("status=%d error=%s\n", apiErr.StatusCode, apiErr.ErrMessage)
    }
}
```

## Testing

```bash
go test ./...
```

Integration tests are skipped unless `OZON_CLIENT_ID` and `OZON_CLIENT_SECRET`
environment variables are set (see `.env.example`).

## License

MIT
