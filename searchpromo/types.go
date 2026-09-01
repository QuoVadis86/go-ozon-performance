package searchpromo

import "github.com/QuoVadis86/go-ozon-performance/transport"

type ActivateAllSkuPromoCampaignResponse struct {
	SelectedBid transport.Uint64 `json:"selectedBid"` // 选定的推广出价。
}

type BatchDisableProductsRequest struct {
	Skus []transport.Uint64 `json:"skus"` // 商品标识符列表。
}

type BatchEnableProductsRequest struct {
	Skus []transport.Uint64 `json:"skus"` // 商品标识符列表。
}

type BatchEnableProductsResponseResponseCell struct {
	Bid    float64          `json:"bid"`    // 投注的含义。
	Error  string           `json:"error"`  // 无法开启商品推广的错误。
	SKU    transport.Uint64 `json:"sku"`    // 推广商品的SKU。
	Update bool             `json:"update"` // 如果商品推广已启用，应为 `true`。
}

type BatchEnableProductsResponse struct {
	Response []BatchEnableProductsResponseResponseCell `json:"response"` // 方法的回复。
}

type DeleteSearchPromoBidsRequestV2 struct {
	SKU []transport.Uint64 `json:"sku"` // 商品标识符列表。
}

type GetProductsRecommendedBidsRequest struct {
	SKU []transport.Uint64 `json:"sku"` // 商品标识符列表。
}

type GetProductsRecommendedBidsResponse struct {
	RecommendedBids map[string]string `json:"recommendedBids"` // 推荐投注列表。
}

type ListSearchPromoProductsRequestV2 struct {
	Page     int64 `json:"page"`     // 页码。分页从1开始。
	PageSize int64 `json:"pageSize"` // 页面尺寸。
}

// 最大投注。
type ListSearchPromoProductsResponseV2Hint struct {
	CampaignId        transport.Uint64 `json:"campaignId"`        // 广告活动标识符。
	OrganisationTitle string           `json:"organisationTitle"` // 组织名称。
}

// 前一轮的投注信息。
type ListSearchPromoProductsResponseV2PreviousBid struct {
	Bid       float64 `json:"bid"`       // 之前的值 `products.bid`。
	BidPrice  string  `json:"bidPrice"`  // 之前的值 `products.bidPrice`。
	UpdatedAt string  `json:"updatedAt"` // 最后更改评分的日期和时间。
}

// 观看信息。
type ListSearchPromoProductsResponseV2Views struct {
	PreviousWeek transport.Uint64 `json:"previousWeek"` // 上周商品的浏览次数。
	ThisWeek     transport.Uint64 `json:"thisWeek"`     // 最近7天内商品的浏览次数。
}

type ListSearchPromoProductsResponseV2Product struct {
	Bid                     float64                                      `json:"bid"`      // 每个订单的佣金（CPO）。单位是商品价格的百分比。
	BidPrice                string                                       `json:"bidPrice"` // 每个订单的佣金（CPO）在卢布。
	Hint                    ListSearchPromoProductsResponseV2Hint        `json:"hint"`
	ImageUrl                string                                       `json:"imageUrl"`               // 图片链接地址。
	IsSearchPromoAvailable  bool                                         `json:"isSearchPromoAvailable"` // `true`，如果商品在“按订单付款”中没有被封锁。
	PreviousBid             ListSearchPromoProductsResponseV2PreviousBid `json:"previousBid"`
	PreviousVisibilityIndex string                                       `json:"previousVisibilityIndex"` // 之前的`visibilityIndex`值。
	Price                   string                                       `json:"price"`                   // 商品价格。
	SearchPromoStatus       bool                                         `json:"searchPromoStatus"`       // `true`，如果已开启在“按订单付款”中的商品推广。
	SKU                     transport.Uint64                             `json:"sku"`                     // 推广商品的SKU。
	SourceSku               string                                       `json:"sourceSku"`               // 卖家货号。
	Title                   string                                       `json:"title"`                   // 商品名称。
	Views                   ListSearchPromoProductsResponseV2Views       `json:"views"`
	VisibilityIndex         string                                       `json:"visibilityIndex"` // 可见性指数。
}

type ListSearchPromoProductsResponseV2 struct {
	Products []ListSearchPromoProductsResponseV2Product `json:"products"` // 商品列表。
	Total    transport.Uint64                           `json:"total"`    // 商品数量。
}

type SetSearchPromoBidsRequestV2Bid struct {
	Bid float64          `json:"bid"` // 每个订单的佣金（CPO），单位是商品价格的百分比。
	SKU transport.Uint64 `json:"sku"` // 商品标识符。
}

type SetSearchPromoBidsRequestV2 struct {
	Bids []SetSearchPromoBidsRequestV2Bid `json:"bids"` // 投注的含义。
}

type SetSearchPromoBidsResponseV2ResponseCell struct {
	Bid    float64          `json:"bid"`    // 投注的含义。
	Error  string           `json:"error"`  // 无法更新评分的错误。
	SKU    transport.Uint64 `json:"sku"`    // 推广商品的SKU。
	Update bool             `json:"update"` // 如果投注已更新，应为 `true`。
}

type SetSearchPromoBidsResponseV2 struct {
	Response []SetSearchPromoBidsResponseV2ResponseCell `json:"response"` // 方法的回复。
}
