package statistics

import "github.com/QuoVadis86/go-ozon-performance/transport"

// GroupBy values
type GroupBy string

const (
	GroupByNOGroupBY    GroupBy = "NO_GROUP_BY"    // 未定义（默认值）
	GroupByDate         GroupBy = "DATE"           // 按日期分组（按天）
	GroupByStartOFWeek  GroupBy = "START_OF_WEEK"  // 按周分组
	GroupByStartOFMonth GroupBy = "START_OF_MONTH" // 按月份分组
)

// StatisticsRequestState values
type StatisticsRequestState string

const (
	StatisticsRequestStateNOTStarted StatisticsRequestState = "NOT_STARTED" // 请求等待执行
	StatisticsRequestStateINProgress StatisticsRequestState = "IN_PROGRESS" // 请求正在进行中
	StatisticsRequestStateError      StatisticsRequestState = "ERROR"       // 请求执行失败
	StatisticsRequestStateOK         StatisticsRequestState = "OK"          // 请求成功完成
)

// ExtstatisticsGroupBy values
type ExtstatisticsGroupBy string

const (
	ExtstatisticsGroupByNOGroupBY    ExtstatisticsGroupBy = "NO_GROUP_BY"    // 未定义（默认值）
	ExtstatisticsGroupByDate         ExtstatisticsGroupBy = "DATE"           // 按日期分组（按天）
	ExtstatisticsGroupByStartOFWeek  ExtstatisticsGroupBy = "START_OF_WEEK"  // 按周分组
	ExtstatisticsGroupByStartOFMonth ExtstatisticsGroupBy = "START_OF_MONTH" // 按月份分组
)

type DownloadStatisticsRequest struct {
	UUID string `url:"UUID"`
}

type GenerateAllSkuPromoOrdersReportRequest struct {
	TimeBoundsFrom string `url:"timeBounds.from"`
	TimeBoundsTo   string `url:"timeBounds.to"`
}

type GenerateAllSkuPromoProductsReportRequest struct {
	TimeBoundsFrom string `url:"timeBounds.from"`
	TimeBoundsTo   string `url:"timeBounds.to"`
}

type GetCampaignDailyStatsRequest struct {
	CampaignIds []transport.Uint64 `url:"campaignIds"`
	DateFrom    string             `url:"dateFrom"`
	DateTo      string             `url:"dateTo"`
}

type GetCampaignExpenseRequest struct {
	CampaignIds []transport.Uint64 `url:"campaignIds"`
	DateFrom    string             `url:"dateFrom"`
	DateTo      string             `url:"dateTo"`
}

type ListReportsExternalRequest struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"pageSize"`
}

type ListReportsRequest struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"pageSize"`
}

type MediaCampaignListRequest struct {
	CampaignIds []transport.Uint64 `url:"campaignIds"`
	DateFrom    string             `url:"dateFrom"`
	DateTo      string             `url:"dateTo"`
	From        string             `url:"from"`
	To          string             `url:"to"`
}

type ProductCampaignListRequest struct {
	CampaignIds []transport.Uint64 `url:"campaignIds"`
	DateFrom    string             `url:"dateFrom"`
	DateTo      string             `url:"dateTo"`
	From        string             `url:"from"`
	To          string             `url:"to"`
}

type SearchPromoProductsSKUStatisticsResponseRow struct {
	AvgCpc      string           `json:"avgCpc"`      // 平均每次点击费用（卢布）。
	CampaignId  transport.Uint64 `json:"campaignId"`  // 广告活动标识符。
	Clicks      transport.Uint64 `json:"clicks"`      // 点击次数。
	CTR         float64          `json:"ctr"`         // 点击次数与展示次数的比例——CTR。
	Date        string           `json:"date"`        // 统计数据对应的日期。
	DateAdded   string           `json:"dateAdded"`   // 广告活动添加日期。
	DRR         float64          `json:"drr"`         // 广告费用份额。
	Expense     string           `json:"expense"`     // 费用（卢布）。
	ModelOrders transport.Uint64 `json:"modelOrders"` // 型号订单数量。
	ModelSales  string           `json:"modelSales"`  // 型号订单金额（卢布）。
	Orders      transport.Uint64 `json:"orders"`      // 订单数量。
	Price       string           `json:"price"`       // 商品价格（卢布）。
	Sales       string           `json:"sales"`       // 订单金额（卢布）。
	SKU         transport.Uint64 `json:"sku"`         // Ozon系统中的商品标识符（SKU）。
	ToCart      transport.Uint64 `json:"toCart"`      // 商品被加入购物车的次数。
	Views       transport.Uint64 `json:"views"`       // 展示次数。
}

type StatisticsReportsListItemCampaign struct {
	ID    transport.Uint64 `json:"id"`    // 广告活动标识符。
	Title string           `json:"title"` // 活动名称。
}

// 原始请求结构。
type StatisticsRequest struct {
	Campaigns []transport.Uint64 `json:"campaigns"` // 需要准备报告的广告活动标识符列表。 报告格式： - CSV — 如果列表中只有一个活动。 - ZIP归档 — 如果列表中有多个广告活动。每个文件对应列表中的一个广告活动。文件名格式为 `<идентификатор кампании>。csv`。
	DateFrom  string             `json:"dateFrom"`  // 报告周期的起始日期格式为YYYY-MM-DD。例如：`2019-02-10`。
	DateTo    string             `json:"dateTo"`    // 报告期的最终日期格式为YYYY-MM-DD。例如：`2019-02-10`。
	From      string             `json:"from"`      // 报告周期的起始日期，以RFC 3339格式。 最长的报告获取时间为62天。
	GroupBy   GroupBy            `json:"groupBy"`
	To        string             `json:"to"` // 报告周期的最终日期，以RFC 3339格式。 最长的报告获取时间为62天。
}

// 请求信息。
// Kind values
type Kind string

const (
	KindStats         Kind = "STATS"          // 项目报告
	KindSearchPhrases Kind = "SEARCH_PHRASES" // 关于搜索词和商品类目的报告
	KindAttribution   Kind = "ATTRIBUTION"    // “按订单付款”的订单报告
	KindVideo         Kind = "VIDEO"          // 视频广告展示报告
)

type StatisticsResponse struct {
	UUID      string                 `json:"UUID"`      // 用于进行检查的请求唯一标识符。
	CreatedAt string                 `json:"createdAt"` // 请求被服务器接收的日期和时间，时区为UTC。
	Error     string                 `json:"error"`     // 出现错误的简要描述。 如果请求执行失败，字段将存在。
	Kind      Kind                   `json:"kind"`      // 请求的报告类型： - `STATS` — 项目报告; - `SEARCH_PHRASES` — 关于搜索词和商品类目的报告； - `ATTRIBUTION`— “按订单付款”的订单报告； - `VIDEO` — 视频广告展示报告。
	Link      string                 `json:"link"`      // 相对链接到CSV格式的报告。 字段存在，如果请求成功。
	Request   StatisticsRequest      `json:"request"`
	State     StatisticsRequestState `json:"state"`
	UpdatedAt string                 `json:"updatedAt"` // 请求状态最后更新时间，时区UTC。
}

type StatisticsReportsListItem struct {
	Campaigns []StatisticsReportsListItemCampaign `json:"campaigns"` // 广告活动列表。
	Meta      StatisticsResponse                  `json:"meta"`
	Name      string                              `json:"name"` // 报告名称。
}

type StatisticsReportsList struct {
	Items []StatisticsReportsListItem `json:"items"` // 报告列表。
	Total transport.Uint64            `json:"total"` // 报告数量。
}

type StatisticsRequestID struct {
	UUID   string `json:"UUID"`   // 已发送请求的唯一标识符。 通过它可以[检查报告生成状态](#operation/VendorStatisticsCheck)和[下载报告](#operation/DownloadStatistics)。
	Vendor bool   `json:"vendor"` // 如果需要外部流量分析报告 — `true`。
}

// 原始请求结构。
type StatisticsVideobannerRequest struct {
	Campaigns []transport.Uint64 `json:"campaigns"` // 需要准备报告的广告活动标识符列表。 报告格式： - CSV — 如果列表中只有一个活动。 - ZIP归档 — 如果列表中有多个广告活动。每个文件对应列表中的一个广告活动。文件名格式为 `<идентификатор кампании>。csv`。
	DateFrom  string             `json:"dateFrom"`  // 报告周期的起始日期格式为YYYY-MM-DD。例如：`2019-02-10`。
	DateTo    string             `json:"dateTo"`    // 报告期的最终日期格式为YYYY-MM-DD。例如：`2019-02-10`。
	GroupBy   GroupBy            `json:"groupBy"`
}

type GenerateAllSkuPromoOrdersReportResponse struct {
	UUID string `json:"UUID"` // 已发送请求的唯一标识符。 通过它可以[检查报告生成状态](#operation/VendorStatisticsCheck)和[下载报告](#operation/DownloadStatistics)。
}

type GenerateAllSkuPromoProductsReportResponse struct {
	UUID string `json:"UUID"` // 已发送请求的唯一标识符。 通过它可以[检查报告生成状态](#operation/VendorStatisticsCheck)和[下载报告](#operation/DownloadStatistics)。
}

type SearchPromoProductsSKUStatisticsRequest struct {
	CampaignIds []transport.Uint64 `json:"campaignIds"` // 广告活动标识符列表。
	DateFrom    string             `json:"dateFrom"`    // 统计周期开始日期，不早于前一天。
	DateTo      string             `json:"dateTo"`      // 统计周期结束日期。
}

type SearchPromoProductsSKUStatisticsResponse struct {
	Rows []SearchPromoProductsSKUStatisticsResponseRow `json:"rows"` // 商品统计数据。
}

type ExtstatisticsStatisticsRequest struct {
	Campaigns []transport.Uint64   `json:"campaigns"` // 广告活动列表。
	DateFrom  string               `json:"dateFrom"`  // 报告周期的起始日期格式为YYYY-MM-DD。例如：`2019-02-10`。 日期不得早于2025年2月1日。
	DateTo    string               `json:"dateTo"`    // 报告期的最终日期格式为YYYY-MM-DD。例如：`2019-02-10`。
	From      string               `json:"from"`      // 报告周期的起始日期，以RFC 3339格式。 日期不得早于2025年2月1日。
	GroupBy   ExtstatisticsGroupBy `json:"groupBy"`
	To        string               `json:"to"` // 报告周期的最终日期，以RFC 3339格式。
}

type ExtstatisticsStatisticsRequestID struct {
	UUID   string `json:"UUID"`   // 已发送请求的唯一标识符。 通过它可以[检查报告生成状态](#operation/VendorStatisticsCheck)和[下载报告](#operation/DownloadStatistics)。
	Vendor bool   `json:"vendor"` // 如果需要外部流量分析报告 — `true`。
}
