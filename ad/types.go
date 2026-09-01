package ad

import "github.com/QuoVadis86/go-ozon-performance/transport"

// Placement values
type Placement string

const (
	PlacementCampaignPlacementInvalid           Placement = "CAMPAIGN_PLACEMENT_INVALID"             // 未定义
	PlacementCampaignPlacementPDP               Placement = "CAMPAIGN_PLACEMENT_PDP"                 // 商品卡片
	PlacementCampaignPlacementSearchANDCategory Placement = "CAMPAIGN_PLACEMENT_SEARCH_AND_CATEGORY" // 搜索和推荐
	PlacementCampaignPlacementTOPPromotion      Placement = "CAMPAIGN_PLACEMENT_TOP_PROMOTION"       // 搜索
)

// DBproductAutopilotStrategy values
type DBproductAutopilotStrategy string

const (
	DBproductAutopilotStrategyMAXClicks      DBproductAutopilotStrategy = "MAX_CLICKS"       // 搜索和推荐的自动策略
	DBproductAutopilotStrategyTOPMAXClicks   DBproductAutopilotStrategy = "TOP_MAX_CLICKS"   // 搜索的自动策略
	DBproductAutopilotStrategyTargetBids     DBproductAutopilotStrategy = "TARGET_BIDS"      // 搜索的平均点击费用
	DBproductAutopilotStrategyNOAutoStrategy DBproductAutopilotStrategy = "NO_AUTO_STRATEGY" // 不使用自动策略
)

// CreateProductCampaignRequestV2ProductCampaignPlacementV2 values
type CreateProductCampaignRequestV2ProductCampaignPlacementV2 string

const (
	CreateProductCampaignRequestV2ProductCampaignPlacementV2PlacementTOPPromotion      CreateProductCampaignRequestV2ProductCampaignPlacementV2 = "PLACEMENT_TOP_PROMOTION"       // 搜索
	CreateProductCampaignRequestV2ProductCampaignPlacementV2PlacementInvalid           CreateProductCampaignRequestV2ProductCampaignPlacementV2 = "PLACEMENT_INVALID"             // 未定义
	CreateProductCampaignRequestV2ProductCampaignPlacementV2PlacementSearchANDCategory CreateProductCampaignRequestV2ProductCampaignPlacementV2 = "PLACEMENT_SEARCH_AND_CATEGORY" // 搜索和推荐
)

// DynamicBudgetType values
type DynamicBudgetType string

const (
	DynamicBudgetTypeDynamicBudgetRequired    DynamicBudgetType = "DYNAMIC_BUDGET_REQUIRED"    // 必填
	DynamicBudgetTypeDynamicBudgetRecommended DynamicBudgetType = "DYNAMIC_BUDGET_RECOMMENDED" // 推荐的
)

// ProductAutopilotStrategyCPC values
type ProductAutopilotStrategyCPC string

const (
	ProductAutopilotStrategyCPCMAXClicks      ProductAutopilotStrategyCPC = "MAX_CLICKS"       // 搜索和推荐的自动策略
	ProductAutopilotStrategyCPCTOPMAXClicks   ProductAutopilotStrategyCPC = "TOP_MAX_CLICKS"   // 搜索的自动策略
	ProductAutopilotStrategyCPCTargetBids     ProductAutopilotStrategyCPC = "TARGET_BIDS"      // 搜索的平均点击费用
	ProductAutopilotStrategyCPCNOAutoStrategy ProductAutopilotStrategyCPC = "NO_AUTO_STRATEGY" // 不使用自动策略
	ProductAutopilotStrategyCPCTOPPromotion   ProductAutopilotStrategyCPC = "TOP_PROMOTION"    // 登上顶端
	ProductAutopilotStrategyCPCTargetCIR      ProductAutopilotStrategyCPC = "TARGET_CIR"       // 目标费用
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

// ProductCampaignPlacement values
type ProductCampaignPlacement string

const (
	ProductCampaignPlacementPlacementInvalid           ProductCampaignPlacement = "PLACEMENT_INVALID"             // 未定义
	ProductCampaignPlacementPlacementPDP               ProductCampaignPlacement = "PLACEMENT_PDP"                 // 商品卡片。仅适用于手动管理的活动
	ProductCampaignPlacementPlacementSearchANDCategory ProductCampaignPlacement = "PLACEMENT_SEARCH_AND_CATEGORY" // 搜索和推荐
	ProductCampaignPlacementPlacementTOPPromotion      ProductCampaignPlacement = "PLACEMENT_TOP_PROMOTION"       // 搜索
)

// 创建活动时的计算参数。
type CalculateDynamicBudgetRequestCreateCampaignScenario struct {
	AutopilotStrategy DBproductAutopilotStrategy `json:"autopilotStrategy"`
	Placement         Placement                  `json:"placement"`
	SkusCount         transport.Uint64           `json:"skusCount"` // 商品数量。
}

// 关于自适应策略的信息。
type MaybeProductAutopilotStrategy struct {
	Strategy DBproductAutopilotStrategy `json:"strategy"`
}

// 更新活动时的计算参数。
type CalculateDynamicBudgetRequestUpdateCampaignScenario struct {
	AddingSkus        []transport.Uint64            `json:"addingSkus"` // 已加入活动的商品。
	AutopilotStrategy MaybeProductAutopilotStrategy `json:"autopilotStrategy"`
	CampaignId        transport.Uint64              `json:"campaignId"` // 广告活动标识符。
}

type Empty struct {
}

// 支付方式： - `CPC` — 点击费用。
type CampaignType string

type CalculateDynamicBudgetRequest struct {
	CreateCampaign CalculateDynamicBudgetRequestCreateCampaignScenario `json:"createCampaign"`
	UpdateCampaign CalculateDynamicBudgetRequestUpdateCampaignScenario `json:"updateCampaign"`
}

type CalculateDynamicBudgetResponse struct {
	CampaignBudget    transport.Uint64  `json:"campaignBudget"` // 广告活动预算。单位是卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
	DynamicBudget     transport.Uint64  `json:"dynamicBudget"`  // 最低预算。
	DynamicBudgetType DynamicBudgetType `json:"dynamicBudgetType"`
	SkuPrice          transport.Uint64  `json:"skuPrice"` // 每向广告活动中添加一件商品时，最低预算增加的金额（以卢布为单位）。
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
	ProductAutopilotStrategyMAXClicks    ProductAutopilotStrategy = "MAX_CLICKS"     // 搜索和推荐的自动策略
	ProductAutopilotStrategyTOPMAXClicks ProductAutopilotStrategy = "TOP_MAX_CLICKS" // 搜索的自动策略
	ProductAutopilotStrategyTargetBids   ProductAutopilotStrategy = "TARGET_BIDS"    // 搜索和推荐的平均点击费用
	ProductAutopilotStrategyTOPPromotion ProductAutopilotStrategy = "TOP_PROMOTION"  // 登上顶端
	ProductAutopilotStrategyTargetCIR    ProductAutopilotStrategy = "TARGET_CIR"     // 目标费用
)

// ProductCampaignMode values
type ProductCampaignMode string

const (
	ProductCampaignModeProductCampaignModeAuto   ProductCampaignMode = "PRODUCT_CAMPAIGN_MODE_AUTO"   // 自动
	ProductCampaignModeProductCampaignModeManual ProductCampaignMode = "PRODUCT_CAMPAIGN_MODE_MANUAL" // 手动
)

type Campaign struct {
	Autopilot                CampaignAutopilotProperties `json:"autopilot"`
	Budget                   transport.Uint64            `json:"budget"`      // 广告活动预算。单位是卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
	CreatedAt                string                      `json:"createdAt"`   // 创建活动的日期格式为RFC3339。
	DailyBudget              transport.Uint64            `json:"dailyBudget"` // 每日广告活动预算。计量单位是千分之一卢布，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
	FromDate                 string                      `json:"fromDate"`    // 广告活动的启动日期。
	ID                       transport.Uint64            `json:"id"`          // 广告活动标识符。
	PaymentType              CampaignType                `json:"paymentType"`
	Placement                []ProductCampaignPlacement  `json:"placement"`                // 广告位置： - `PLACEMENT_INVALID` — 未定义； - `PLACEMENT_PDP` — 商品卡片； - `PLACEMENT_SEARCH_AND_CATEGORY` — 搜索和类目（模版）； - `PLACEMENT_TOP_PROMOTION` — 输出到顶部。
	ProductAutopilotStrategy ProductAutopilotStrategy    `json:"productAutopilotStrategy"` // 广告活动当前使用的策略： - `MAX_CLICKS` — 搜索和推荐的自动策略； - `TOP_MAX_CLICKS` — 搜索的自动策略； - `TARGET_BIDS` — 搜索和推荐的平均点击费用； - `TOP_PROMOTION` — 登上顶端； - `TARGET_CIR` — 目标费用。
	ProductCampaignMode      ProductCampaignMode         `json:"productCampaignMode"`      // 创建和管理商品广告活动模式： - `PRODUCT_CAMPAIGN_MODE_AUTO` — 自动; - `PRODUCT_CAMPAIGN_MODE_MANUAL` — 手动。
	State                    CampaignState               `json:"state"`
	Title                    string                      `json:"title"`        // 活动名称。
	ToDate                   string                      `json:"toDate"`       // 广告活动结束日期。
	UpdatedAt                string                      `json:"updatedAt"`    // 更新广告的日期，以RFC3339格式。
	WeeklyBudget             transport.Uint64            `json:"weeklyBudget"` // 每周广告活动预算。单位是卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。
}

type CampaignID struct {
	CampaignId transport.Uint64 `json:"campaignId"` // 广告活动标识符。
}

type CreateProductCampaignRequestV2CPC struct {
	DailyBudget              transport.Uint64                                         `json:"dailyBudget"` // 广告活动每日预算限制。单位是俄罗斯卢布的百万分之一，四舍五入到戈比。 例如，参数值 `1 000 000` 等于1卢布。 如果参数未填写，日预算不受限制。创建广告活动后，无法将日预算更改为周预算或反之。
	FromDate                 string                                                   `json:"fromDate"`    // 广告活动开始日期（莫斯科时间）。 如果参数未填写，则启动日期为当前天的开始。 如果不需要审核，解封操作会在活动激活后立即开始。
	Placement                CreateProductCampaignRequestV2ProductCampaignPlacementV2 `json:"placement"`
	ProductAutopilotStrategy ProductAutopilotStrategyCPC                              `json:"productAutopilotStrategy"`
	Title                    string                                                   `json:"title"`        // 广告活动名称。
	ToDate                   string                                                   `json:"toDate"`       // 按莫斯科时间的广告活动结束日期。 对于在自动模式下创建的点击付费广告活动，不考虑此参数。
	WeeklyBudget             transport.Uint64                                         `json:"weeklyBudget"` // 广告活动周预算限制。单位是卢布的百万分之一，四舍五入到戈比。 例如，参数值 `1 000 000` 等于1卢布。 如果参数未填写，周预算不受限制。创建广告活动后，无法将周预算更改为日预算或反之。
}

type PatchProductCampaignRequest struct {
	Autopilot    map[string]any   `json:"autopilot"`    // 活动信息。 如果参数 `productAutopilotStrategy` 开启了自动策略，则为必填项。
	Budget       transport.Uint64 `json:"budget"`       // 广告活动总预算限制。单位是卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。 仅适用于品牌和代理机构的自动化活动。 在其他组织中无法安装新的总预算限制。如果已为活动设置预算，可以： - 移除限制：将此参数传为 `0`。 - 不要更改公司预算。当预算用完时，可以删除预算限制或创建一个无限预算的新活动。
	DailyBudget  transport.Uint64 `json:"dailyBudget"`  // 广告活动每日预算限制。单位是俄罗斯卢布的百万分之一，四舍五入到分。例如，参数值 `1 000 000` 等于1卢布。 如果参数未填写，日预算不受限制。创建广告活动后，无法将日预算更改为周预算或反之。
	FromDate     string           `json:"fromDate"`     // 广告活动开始日期（莫斯科时间）。 不能早于当前日期。
	ToDate       string           `json:"toDate"`       // 按莫斯科时间的广告活动结束日期。 不能早于开始日期。
	WeeklyBudget transport.Uint64 `json:"weeklyBudget"` // 广告活动周预算限制。单位是卢布的百万分之一，四舍五入到戈比。 例如，参数值 `1 000 000` 等于1卢布。 如果参数未填写，周预算不受限制。创建广告活动后，无法将周预算更改为日预算或反之。
}
