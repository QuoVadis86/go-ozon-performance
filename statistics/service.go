package statistics

import (
	"context"
	"strings"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

type Service struct{ Client *transport.Client }

// 订单报告
func (s *Service) AttributionSubmitRequest(ctx context.Context, req *StatisticsRequest) (*StatisticsRequestID, error) {
	path := "/api/client/statistics/attribution"
	var resp StatisticsRequestID
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 获取报告
func (s *Service) DownloadStatistics(ctx context.Context, req *DownloadStatisticsRequest) ([]byte, error) {
	path := "/api/client/statistics/report"
	return s.Client.GetRaw(ctx, path, req)
}

// 获取按订单付费订单报告——所有商品
func (s *Service) GenerateAllSkuPromoOrdersReport(ctx context.Context, req *GenerateAllSkuPromoOrdersReportRequest) (*GenerateAllSkuPromoOrdersReportResponse, error) {
	path := "/api/client/statistics/all_sku_promo/orders/generate"
	var resp GenerateAllSkuPromoOrdersReportResponse
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 获取按订单付费商品报告——所有商品
func (s *Service) GenerateAllSkuPromoProductsReport(ctx context.Context, req *GenerateAllSkuPromoProductsReportRequest) (*GenerateAllSkuPromoProductsReportResponse, error) {
	path := "/api/client/statistics/all_sku_promo/products/generate"
	var resp GenerateAllSkuPromoProductsReportResponse
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 每日广告统计
func (s *Service) GetCampaignDailyStats(ctx context.Context, req *GetCampaignDailyStatsRequest) ([]byte, error) {
	path := "/api/client/statistics/daily"
	return s.Client.GetRaw(ctx, path, req)
}

// 广告活动费用统计
func (s *Service) GetCampaignExpense(ctx context.Context, req *GetCampaignExpenseRequest) ([]byte, error) {
	path := "/api/client/statistics/expense"
	return s.Client.GetRaw(ctx, path, req)
}

// 通过界面生成的报告列表
func (s *Service) ListReports(ctx context.Context, req *ListReportsRequest) (*StatisticsReportsList, error) {
	path := "/api/client/statistics/list"
	var resp StatisticsReportsList
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 通过API生成的报告列表
func (s *Service) ListReportsExternal(ctx context.Context, req *ListReportsExternalRequest) (*StatisticsReportsList, error) {
	path := "/api/client/statistics/externallist"
	var resp StatisticsReportsList
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 媒体活动统计
func (s *Service) MediaCampaignList(ctx context.Context, req *MediaCampaignListRequest) ([]byte, error) {
	path := "/api/client/statistics/campaign/media"
	return s.Client.GetRaw(ctx, path, req)
}

// 按点击付费和特别投放活动统计
func (s *Service) ProductCampaignList(ctx context.Context, req *ProductCampaignListRequest) ([]byte, error) {
	path := "/api/client/statistics/campaign/product"
	return s.Client.GetRaw(ctx, path, req)
}

// 获取按订单付费订单报告——所选商品
func (s *Service) SearchPromoOrdersReportSubmitRequest(ctx context.Context) (*StatisticsRequestID, error) {
	path := "/api/client/statistic/orders/generate"
	var resp StatisticsRequestID
	if err := s.Client.Post(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 获取按订单付费商品报告——所选商品
func (s *Service) SearchPromoProductsReportSubmitRequest(ctx context.Context) (*StatisticsRequestID, error) {
	path := "/api/client/statistic/products/generate"
	var resp StatisticsRequestID
	if err := s.Client.Post(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 获取按点击付费广告活动中的商品统计数据
func (s *Service) SearchPromoProductsSKUStatistics2(ctx context.Context, req *SearchPromoProductsSKUStatisticsRequest) (*SearchPromoProductsSKUStatisticsResponse, error) {
	path := "/api/client/statistics/products/sku"
	var resp SearchPromoProductsSKUStatisticsResponse
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 报告状态
func (s *Service) StatisticsCheck(ctx context.Context, uUID string) (*StatisticsResponse, error) {
	path := "/api/client/statistics/{UUID}"
	path = strings.Replace(path, "{UUID}", uUID, 1)
	var resp StatisticsResponse
	if err := s.Client.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 搜索查询报告
func (s *Service) StatisticsPhrasesSubmitRequest2(ctx context.Context, req *ExtstatisticsStatisticsRequest) (*ExtstatisticsStatisticsRequestID, error) {
	path := "/api/client/statistics/phrases"
	var resp ExtstatisticsStatisticsRequestID
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 广告活动统计
func (s *Service) SubmitRequest(ctx context.Context, req *StatisticsRequest) (*StatisticsRequestID, error) {
	path := "/api/client/statistics"
	var resp StatisticsRequestID
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 视频横幅展示统计
func (s *Service) VideoCampaignsSubmitRequest(ctx context.Context, req *StatisticsVideobannerRequest) (*StatisticsRequestID, error) {
	path := "/api/client/statistics/video"
	var resp StatisticsRequestID
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
