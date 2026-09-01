package campaign

import "github.com/QuoVadis86/go-ozon-performance/transport"

// CampaignTypeInList values
type CampaignTypeInList string

const (
	CampaignTypeInListCPC                 CampaignTypeInList = "CPC"                   // 点击费用
	CampaignTypeInListCPM                 CampaignTypeInList = "CPM"                   // 按展示计价
	CampaignTypeInListCPO                 CampaignTypeInList = "CPO"                   // 为订单
	CampaignTypeInListCampaignTypeInvalid CampaignTypeInList = "CAMPAIGN_TYPE_INVALID" // 未定义（默认值）
)

// CampaignTypeRate values
type CampaignTypeRate string

const (
	CampaignTypeRateCPO    CampaignTypeRate = "CPO"     // “按订单付款” 类型的活动
	CampaignTypeRateCPC    CampaignTypeRate = "CPC"     // 搜索和推荐
	CampaignTypeRateCPCTOP CampaignTypeRate = "CPC_TOP" // 搜索
)

// MarketplaceID values
type MarketplaceID string

const (
	MarketplaceIDMarketplaceIDRU MarketplaceID = "MARKETPLACE_ID_RU" // 俄罗斯
	MarketplaceIDMarketplaceIDKZ MarketplaceID = "MARKETPLACE_ID_KZ" // 哈萨克斯坦
	MarketplaceIDMarketplaceIDBY MarketplaceID = "MARKETPLACE_ID_BY" // 白俄罗斯
)

// CampaignState values
type CampaignState string

const (
	CampaignStateCampaignStateRunning              CampaignState = "CAMPAIGN_STATE_RUNNING"                // 活跃的活动
	CampaignStateCampaignStatePlanned              CampaignState = "CAMPAIGN_STATE_PLANNED"                // 计划中的活动，时间还没有到
	CampaignStateCampaignStateStopped              CampaignState = "CAMPAIGN_STATE_STOPPED"                // 由于预算不足而暂停的活动
	CampaignStateCampaignStateInactive             CampaignState = "CAMPAIGN_STATE_INACTIVE"               // 由于拥有者停止而停止的活动
	CampaignStateCampaignStateArchived             CampaignState = "CAMPAIGN_STATE_ARCHIVED"               // 归档的活动
	CampaignStateCampaignStateModerationDraft      CampaignState = "CAMPAIGN_STATE_MODERATION_DRAFT"       // 已编辑的未发送到审核的活动
	CampaignStateCampaignStateModerationINProgress CampaignState = "CAMPAIGN_STATE_MODERATION_IN_PROGRESS" // 已提交审核的广告活动
	CampaignStateCampaignStateModerationFailed     CampaignState = "CAMPAIGN_STATE_MODERATION_FAILED"      // 未通过审核的广告活动
	CampaignStateCampaignStateFinished             CampaignState = "CAMPAIGN_STATE_FINISHED"               // 该活动已结束，结束日期已过去，无法更改此活动，可以只能克隆或创建新的一份
	CampaignStateCampaignStateUnknown              CampaignState = "CAMPAIGN_STATE_UNKNOWN"                // 未定义（默认值）
)

// ProductCampaignInListPlacement values
type ProductCampaignInListPlacement string

const (
	ProductCampaignInListPlacementPlacementInvalid           ProductCampaignInListPlacement = "PLACEMENT_INVALID"             // 未定义
	ProductCampaignInListPlacementPlacementPDP               ProductCampaignInListPlacement = "PLACEMENT_PDP"                 // 商品卡片。仅适用于手动管理的活动
	ProductCampaignInListPlacementPlacementSearchANDCategory ProductCampaignInListPlacement = "PLACEMENT_SEARCH_AND_CATEGORY" // 搜索和推荐
	ProductCampaignInListPlacementPlacementTOPPromotion      ProductCampaignInListPlacement = "PLACEMENT_TOP_PROMOTION"       // 搜索
	ProductCampaignInListPlacementPlacementOvertop           ProductCampaignInListPlacement = "PLACEMENT_OVERTOP"             // 搜索和首页（特别投放）
	ProductCampaignInListPlacementPlacementTakeover          ProductCampaignInListPlacement = "PLACEMENT_TAKEOVER"            // 同时在前4个模块中展示商品
	ProductCampaignInListPlacementPlacementTileMain          ProductCampaignInListPlacement = "PLACEMENT_TILE_MAIN"           // 首页，原生横幅
)

type BidBySKUResponseBidSKU struct {
	Bid float64          `json:"bid"` // 出价金额。
	SKU transport.Uint64 `json:"sku"` // 商品标识符：Ozon ID 或 SKU。
}

type ListCampaignsRequest struct {
	AdvObjectType string             `url:"advObjectType"`
	CampaignIds   []transport.Uint64 `url:"campaignIds"`
	Page          int64              `url:"page"`
	PageSize      int64              `url:"pageSize"`
	State         CampaignState      `url:"state"`
}

type BidBySKURequest struct {
	MarketplaceId MarketplaceID      `json:"marketplaceId"`
	PaymentType   CampaignTypeRate   `json:"paymentType"`
	SKU           []transport.Uint64 `json:"sku"` // 商品标识符：Ozon ID 或 SKU
}

type BidBySKUResponse struct {
	MinBids []BidBySKUResponseBidSKU `json:"minBids"` // 关于最低出价的信息。
}

// 活动信息。如果在`productAutopilotStrategy`参数中启用了策略，则返回该信息。
// SkuAddMode values
type SkuAddMode string

const (
	SkuAddModeProductCampaignSKUADDModeUnknown SkuAddMode = "PRODUCT_CAMPAIGN_SKU_ADD_MODE_UNKNOWN" // 商品添加到活动的策略未设置
	SkuAddModeProductCampaignSKUADDModeManual  SkuAddMode = "PRODUCT_CAMPAIGN_SKU_ADD_MODE_MANUAL"  // 只能手动将商品添加到活动中
	SkuAddModeProductCampaignSKUADDModeAuto    SkuAddMode = "PRODUCT_CAMPAIGN_SKU_ADD_MODE_AUTO"    // 可以自动将类目 `categoryId` 的商品添加到活动中
)

type CampaignAutopilotProperties struct {
	CategoryId transport.Uint64 `json:"categoryId"` // 活动商品类目标识符。
	SkuAddMode SkuAddMode       `json:"skuAddMode"` // 是否允许将`categoryId`中指定类目的商品自动添加到采用`MAX_VIEWS`策略的广告活动中： - `PRODUCT_CAMPAIGN_SKU_ADD_MODE_UNKNOWN` — 商品添加到活动的策略未设置。 - `PRODUCT_CAMPAIGN_SKU_ADD_MODE_MANUAL` — 只能手动将商品添加到活动中。 - `PRODUCT_CAMPAIGN_SKU_ADD_MODE_AUTO` — 可以自动将类目 `categoryId` 的商品添加到活动中。 如果许可类型为空或已传递为 `PRODUCT_CAMPAIGN_SKU_ADD_MODE_UNKNOWN`，只能手动将商品添加到活动中。 对于未采用`MAX_VIEWS`策略的广告活动，返回`PRODUCT_CAMPAIGN_SKU_ADD_MODE_UNKNOWN`。
}

// ProductAutopilotStrategy values
type ProductAutopilotStrategy string

const (
	ProductAutopilotStrategyMAXViews       ProductAutopilotStrategy = "MAX_VIEWS"        // 最大显示次数
	ProductAutopilotStrategyMAXClicks      ProductAutopilotStrategy = "MAX_CLICKS"       // 搜索和推荐的自动策略
	ProductAutopilotStrategyTOPMAXClicks   ProductAutopilotStrategy = "TOP_MAX_CLICKS"   // 搜索的自动策略
	ProductAutopilotStrategyNOAutoStrategy ProductAutopilotStrategy = "NO_AUTO_STRATEGY" // 不使用自动策略
	ProductAutopilotStrategyTakeover       ProductAutopilotStrategy = "TAKEOVER"         // 搜索特别投放
	ProductAutopilotStrategyTOPPromotion   ProductAutopilotStrategy = "TOP_PROMOTION"    // 登上顶端
)

// ProductCampaignMode values
type ProductCampaignMode string

const (
	ProductCampaignModeProductCampaignModeAuto   ProductCampaignMode = "PRODUCT_CAMPAIGN_MODE_AUTO"   // 自动
	ProductCampaignModeProductCampaignModeManual ProductCampaignMode = "PRODUCT_CAMPAIGN_MODE_MANUAL" // 手动
)

type CampaignInList struct {
	AdvObjectType            string                           `json:"advObjectType"` // 广告活动类型： - `SKU` — 按点击付费； - `SEARCH_PROMO` — 按订单付款。
	Autopilot                CampaignAutopilotProperties      `json:"autopilot"`
	Budget                   transport.Uint64                 `json:"budget"`      // 广告活动预算。单位是卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
	CreatedAt                string                           `json:"createdAt"`   // 创建活动的日期格式为RFC3339。
	DailyBudget              transport.Uint64                 `json:"dailyBudget"` // 每日广告活动预算。计量单位是千分之一卢布，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
	FromDate                 string                           `json:"fromDate"`    // 广告活动的启动日期。
	ID                       transport.Uint64                 `json:"id"`          // 广告活动标识符。
	PaymentType              CampaignTypeInList               `json:"paymentType"`
	Placement                []ProductCampaignInListPlacement `json:"placement"`                // 广告位置： - `PLACEMENT_INVALID` — 未定义； - `PLACEMENT_PDP` — 商品卡片； - `PLACEMENT_SEARCH_AND_CATEGORY` — 搜索和类目（模版）。 - `PLACEMENT_TOP_PROMOTION` — 输出到顶部。 - `PLACEMENT_OVERTOP` — 搜索和首页（特别投放）。 - `PLACEMENT_TAKEOVER` — 同时在前4个模块中展示商品。 - `PLACEMENT_TILE_MAIN` — 首页，原生横幅。
	ProductAutopilotStrategy ProductAutopilotStrategy         `json:"productAutopilotStrategy"` // 广告活动当前使用的策略： - `MAX_VIEWS` — 最大显示次数； - `MAX_CLICKS` — 搜索和推荐的自动策略； - `TOP_MAX_CLICKS` — 搜索的自动策略； - `NO_AUTO_STRATEGY` — 不使用自动策略； - `TAKEOVER` — 搜索特别投放； - `TOP_PROMOTION` — 登上顶端。
	ProductCampaignMode      ProductCampaignMode              `json:"productCampaignMode"`      // 创建和管理商品广告活动模式： - `PRODUCT_CAMPAIGN_MODE_AUTO` — 自动; - `PRODUCT_CAMPAIGN_MODE_MANUAL` — 手动。
	State                    CampaignState                    `json:"state"`
	Title                    string                           `json:"title"`        // 活动名称。
	ToDate                   string                           `json:"toDate"`       // 广告活动结束日期。
	UpdatedAt                string                           `json:"updatedAt"`    // 更新广告的日期，以RFC3339格式。
	WeeklyBudget             transport.Uint64                 `json:"weeklyBudget"` // 每周广告活动预算。单位是卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
}

type CampaignObject struct {
	ID transport.Uint64 `json:"id"` // 广告对象标识符： - SKU — 用于在赞助商区和类目中的商品广告； - 数字标识符 — 用于横幅广告活动。
}

type CampaignObjectsList struct {
	List []CampaignObject `json:"list"` // 广告对象标识符列表。
}

type CampaignsList struct {
	List []CampaignInList `json:"list"` // 广告活动列表。
}

type CategoriesLimits struct {
	Bid      float64 `json:"bid"`      // 在卢布的利率。
	Category string  `json:"category"` // 商品类别。
}

// ObjectType values
type ObjectType string

const (
	ObjectTypeSKU         ObjectType = "SKU"          // 按点击付费
	ObjectTypeSearchPromo ObjectType = "SEARCH_PROMO" // 按订单付款
)

// PaymentMethod values
type PaymentMethod string

const (
	PaymentMethodCPO PaymentMethod = "CPO" // 订单付款
	PaymentMethodCPC PaymentMethod = "CPC" // 点击付费
	PaymentMethodCPM PaymentMethod = "CPM" // 每 1000 次展示支付
)

// Placement values
type Placement string

const (
	PlacementCampaignPlacementSearchANDCategory Placement = "CAMPAIGN_PLACEMENT_SEARCH_AND_CATEGORY" // 搜索和推荐
	PlacementCampaignPlacementTOPPromotion      Placement = "CAMPAIGN_PLACEMENT_TOP_PROMOTION"       // 搜索
	PlacementCampaignPlacementOvertop           Placement = "CAMPAIGN_PLACEMENT_OVERTOP"             // 特别投放
)

type LimitsData struct {
	Categories    []CategoriesLimits `json:"categories"`    // 二级类目最低费率。如果类目不在列表中，则应用工具的最低费率。
	MaxBid        float64            `json:"maxBid"`        // 最大投注金额（俄罗斯卢布）。
	MinBid        float64            `json:"minBid"`        // 最低押金（以卢布为单位）。
	ObjectType    ObjectType         `json:"objectType"`    // 工具类型： - `SKU` — 按点击付费， - `SEARCH_PROMO` — 按订单付款。
	PaymentMethod PaymentMethod      `json:"paymentMethod"` // 付款方式： - `CPO` — 订单付款， - `CPC` — 点击付费， - `CPM` — 每 1000 次展示支付。
	Placement     Placement          `json:"placement"`     // 广告类型： - `CAMPAIGN_PLACEMENT_SEARCH_AND_CATEGORY` — 搜索和推荐； - `CAMPAIGN_PLACEMENT_TOP_PROMOTION` — 搜索； - `CAMPAIGN_PLACEMENT_OVERTOP` — 特别投放。
}

type ListLimitsResponse struct {
	Limits []LimitsData `json:"limits"` // 推广工具的最低和最高投注。
}

type ListProductsWithBonusesResponse struct {
	Skus []transport.Uint64 `json:"skus"` // 商品SKU列表。
}
