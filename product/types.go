package product

import "github.com/QuoVadis86/go-ozon-performance/transport"

type GetProductsCompetitiveBidsRequest struct {
	Skus []string `url:"skus"`
}

type GetProductsCompetitiveBidsResponseBids struct {
	Bid transport.Uint64 `json:"bid"` // 赔偿率。
	SKU string           `json:"sku"` // 商品SKU。
}

type GetProductsV2Request struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"pageSize"`
}

type AddProductsRequestProduct struct {
	Bid         transport.Uint64 `json:"bid"`         // <aside class="notice"> 仅适用于“平均每次点击费用”策略。</aside> 单次点击出价——CPC。
	SKU         transport.Uint64 `json:"sku"`         // 促销商品的SKU。 必填字段。
	TargetCir   float64          `json:"targetCir"`   // <aside class="notice"> 仅适用于“目标费用”策略。</aside> 目标广告费用份额，单位为百分比。
	TopPosition transport.Uint64 `json:"topPosition"` // <aside class="notice"> 仅适用于采用“登上顶端”策略的广告活动。</aside> 搜索结果中的位置范围： - `最佳-4`， - `最佳-12`， - `最佳-20`， - `最佳-30`.
}

type AddProductsRequest struct {
	Bids []AddProductsRequestProduct `json:"bids"` // 点击佣金。
}

type DeleteProductsRequest struct {
	SKU []transport.Uint64 `json:"sku"` // 推广商品的SKU。
}

type GetProductsCompetitiveBidsResponse struct {
	Bids       []GetProductsCompetitiveBidsResponseBids `json:"bids"`       // 投注列表。
	CampaignId string                                   `json:"campaignId"` // 广告活动标识符。
}

type GetProductsResponseV2Product struct {
	Bid         transport.Uint64 `json:"bid"`         // <aside class="notice"> 仅适用于“平均每次点击费用”策略。</aside> 单次点击出价——CPC。
	SKU         transport.Uint64 `json:"sku"`         // 商品标识：Ozon ID 或 SKU。
	TargetCir   float64          `json:"targetCir"`   // <aside class="notice"> 仅适用于“目标费用”策略。</aside> 目标广告费用份额，单位为百分比。
	Title       string           `json:"title"`       // 商品名称。
	TopPosition transport.Uint64 `json:"topPosition"` // <aside class="notice"> 仅适用于采用“登上顶端”策略的广告活动。</aside> 搜索结果中的位置范围： - `最佳-4`， - `最佳-12`， - `最佳-20`， - `最佳-30`.
}

type GetProductsResponseV2 struct {
	Products []GetProductsResponseV2Product `json:"products"` // 活动商品列表。
}

type UpdateProductsRequestProduct struct {
	Bid         transport.Uint64 `json:"bid"`         // <aside class="notice"> 仅适用于“平均每次点击费用”策略。</aside> 单次点击出价——CPC。
	SKU         transport.Uint64 `json:"sku"`         // 促销商品的SKU。
	TargetCir   float64          `json:"targetCir"`   // <aside class="notice"> 仅适用于“目标费用”策略。</aside> 目标广告费用份额，单位为百分比。
	TopPosition transport.Uint64 `json:"topPosition"` // <aside class="notice"> 仅适用于采用“登上顶端”策略的广告活动。</aside> 搜索结果中的位置范围： - `最佳-4`， - `最佳-12`， - `最佳-20`， - `最佳-30`.
}

type UpdateProductsRequest struct {
	Bids []UpdateProductsRequestProduct `json:"bids"` // 点击佣金。
}
