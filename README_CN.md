# Ozon Performance API Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?logo=go)](https://golang.org/dl/)
[![Go Reference](https://pkg.go.dev/badge/github.com/QuoVadis86/go-ozon-performance.svg)](https://pkg.go.dev/github.com/QuoVadis86/go-ozon-performance)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English Documentation](README.md)

**go-ozon-performance** 是 [Ozon Performance API](https://docs.ozon.ru/api/performance/) 的 Go 语言客户端库——俄罗斯电商平台 Ozon 广告平台的 API，覆盖活动、统计、报表、商品出价与外部流量分析。  
覆盖 **44 个 API 方法**，**6 个服务模块**，**115 个生成类型**。

## 特性

- 完整覆盖 Ozon Performance API v2（OpenAPI 3.0）
- 自动 OAuth2 令牌管理——自动获取、缓存、透明刷新 Bearer 令牌
- 基于 `swagger.json` 生成的强类型请求/响应结构
- `transport.Uint64` 同时兼容 JSON 数字与字符串两种线格式
- 异步报告流程支持（提交 → 查询 → 下载 CSV/ZIP）
- 带文档说明的枚举类型与描述注释
- 支持 Context 取消与超时控制
- 零外部依赖（仅使用 Go 标准库）
- 并发安全的 HTTP 连接池

## 安装

```bash
go get github.com/QuoVadis86/go-ozon-performance
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    ozon "github.com/QuoVadis86/go-ozon-performance"
)

func main() {
    // 使用凭证创建客户端（个人中心 → 设置 → API-密钥 → Performance API）
    client := ozon.NewClient("your-client-id", "your-client-secret", nil)

    // 获取广告活动列表
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

## API 覆盖

| 服务 | 方法数 | 说明 |
|---------|--------|------|
| [Campaign](./campaign/service.go) | 5 | 活动列表、推广对象、投注限额、按 SKU 最低出价 |
| [Statistics](./statistics/service.go) | 17 | 活动统计报告、搜索词统计、SKU 统计 |
| [Ad](./ad/service.go) | 5 | 活动创建、启动、暂停、修改、动态预算 |
| [Product](./product/service.go) | 5 | 按点击付费活动中的商品出价管理 |
| [Search-Promo](./searchpromo/service.go) | 8 | 按订单付款活动、推荐出价 |
| [Vendor](./vendor/service.go) | 4 | 外部流量分析报告 |

## 示例

### 获取活动列表

```go
import "github.com/QuoVadis86/go-ozon-performance/campaign"

resp, err := client.Campaign.ListCampaigns(ctx, &campaign.ListCampaignsRequest{
    AdvObjectType: "SKU",
    State:         campaign.StateCampaignStateRunning, // 仅运行中的活动
    Page:          1,
    PageSize:      100,
})
```

### 提交异步统计报告

```go
import (
    "github.com/QuoVadis86/go-ozon-performance/statistics"
    "github.com/QuoVadis86/go-ozon-performance/transport"
)

// 1. 提交报告请求——返回 UUID
sub, err := client.Statistics.SubmitRequest(ctx, &statistics.StatisticsRequest{
    Campaigns: []transport.Uint64{48852},
    DateFrom:  "2026-08-01",
    DateTo:    "2026-08-31",
    GroupBy:   statistics.GroupByDate,
})

// 2. 轮询直到生成完成
for {
    st, err := client.Statistics.StatisticsCheck(ctx, sub.UUID)
    if err != nil {
        log.Fatal(err)
    }
    if st.State == statistics.StatisticsRequestStateOK {
        // 3. 下载 CSV 报告
        data, err := client.Statistics.DownloadStatistics(ctx, &statistics.DownloadStatisticsRequest{
            UUID: sub.UUID,
        })
        // 保存 data 到文件...
        break
    }
    time.Sleep(5 * time.Second)
}
```

### 管理商品出价

```go
import (
    "github.com/QuoVadis86/go-ozon-performance/product"
    "github.com/QuoVadis86/go-ozon-performance/transport"
)

// 向按点击付费活动添加商品
err := client.Product.AddProducts(ctx, campaignID, &product.AddProductsRequest{
    Bids: []product.AddProductsRequestProduct{
        {
            SKU: transport.Uint64(48852),
            Bid: transport.Uint64(100000), // 单位为卢布百万分之一
        },
    },
})
```

## 令牌管理

令牌获取是自动的。`POST /api/client/token` 只会触发一次，在并发调用中缓存，
并在过期前自动刷新。收到 `401` 响应会强制刷新并重试一次。进阶场景可通过
`TokenStore` 在进程重启间持久化令牌：

```go
client := ozon.NewClient(clientID, secret, &ozon.ClientOptions{
    TokenStore: myPersistentStore, // 实现 transport.TokenStore 接口
})
```

## 错误处理

所有 API 错误响应返回 `*transport.APIError`，包含状态码与错误信息。用
`errors.As` 检查：

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

## 测试

```bash
go test ./...
```

集成测试在未设置 `OZON_CLIENT_ID` 和 `OZON_CLIENT_SECRET` 环境变量时自动跳过
（参见 `.env.example`）。

## 许可证

MIT
