package vendor

import "github.com/QuoVadis86/go-ozon-performance/transport"

type GetVendorTagRequest struct {
	OrgId transport.Uint64 `url:"orgId"`
}

type VendorStatisticsCheckRequest struct {
	Vendor bool `url:"vendor"`
}

type VendorStatisticsListReportsRequest struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"pageSize"`
}

// 请求生成报告的方法。
// Type values
type Type string

const (
	TypeTrafficSources Type = "TRAFFIC_SOURCES" // 交通来源报告
	TypeOrders         Type = "ORDERS"          // 订单报告
)

type VendorStatisticsRequest struct {
	DateFrom string `json:"dateFrom"` // 报告期开始。 日期不得早于2022年1月1日。
	DateTo   string `json:"dateTo"`   // 报告期结束。 日期必须晚于 `dateFrom` 不超过3个月。如果日期晚于3个月， 报告将在指定在 `dateFrom` 中的日期后3个月内生成。
	Type     Type   `json:"type"`     // 报告类型： - `TRAFFIC_SOURCES` — 交通来源报告。 - `ORDERS` — 订单报告。
}

// 报告信息。
// State values
type State string

const (
	StateNOTStarted State = "NOT_STARTED" // 仍未开始
	StateINProgress State = "IN_PROGRESS" // 在进行中
	StateOK         State = "OK"          // 已取消
	StateError      State = "ERROR"       // 报告已生成
	StateTimeout    State = "TIMEOUT"     // 发生错误
	StateCancel     State = "CANCEL"      // 等待时间已到
)

type VendorStatisticsResponse struct {
	UUID      string                  `json:"UUID"`      // 已发送请求的唯一标识符，可以用于后续检查请求执行状态并下载报告。
	CreatedAt string                  `json:"createdAt"` // 报告生成日期。
	Error     string                  `json:"error"`     // 错误描述。
	Link      string                  `json:"link"`      // 报告链接。
	Request   VendorStatisticsRequest `json:"request"`
	State     State                   `json:"state"`     // 报告生成状态： - `NOT_STARTED` — 仍未开始。 - `IN_PROGRESS` — 在进行中。 - `CANCEL` — 已取消。 - `OK` — 报告已生成。 - `ERROR` — 发生错误。 - `TIMEOUT` — 等待时间已到。
	UpdatedAt string                  `json:"updatedAt"` // 状态变更时间。
}

type VendorStatisticsReportsListItem struct {
	Meta VendorStatisticsResponse `json:"meta"`
	Name string                   `json:"name"` // 报告名称。
}

type GetVendorTagResponse struct {
	Tag string `json:"tag"` // 组织的前缀用于UTM标记。
}

type StatisticsRequestID struct {
	UUID   string `json:"UUID"`   // 已发送请求的唯一标识符。 通过它可以[检查报告生成状态](#operation/VendorStatisticsCheck)和[下载报告](#operation/DownloadStatistics)。
	Vendor bool   `json:"vendor"` // 如果需要外部流量分析报告 — `true`。
}

type VendorStatisticsReportsList struct {
	Items []VendorStatisticsReportsListItem `json:"items"` // 报告信息。
	Total transport.Uint64                  `json:"total"` // 报告数量。
}
