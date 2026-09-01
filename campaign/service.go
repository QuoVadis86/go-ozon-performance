package campaign

import (
	"context"
	"strings"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

type Service struct{ Client *transport.Client }

// 按 SKU 获取商品最低出价
func (s *Service) BidBySKU(ctx context.Context, req *BidBySKURequest) (*BidBySKUResponse, error) {
	path := "/api/min/sku"
	var resp BidBySKUResponse
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 推广工具的投注限额
func (s *Service) GetLimitsList(ctx context.Context) (*ListLimitsResponse, error) {
	path := "/api/client/limits/list"
	var resp ListLimitsResponse
	if err := s.Client.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 活动中的推广对象列表
func (s *Service) ListCampaignObjects(ctx context.Context, campaignId transport.Uint64) (*CampaignObjectsList, error) {
	path := "/api/client/campaign/{campaignId}/objects"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	var resp CampaignObjectsList
	if err := s.Client.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 广告活动列表
func (s *Service) ListCampaigns(ctx context.Context, req *ListCampaignsRequest) (*CampaignsList, error) {
	path := "/api/client/campaign"
	var resp CampaignsList
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 带有积分的商品列表
func (s *Service) ListProductsWithBonuses2(ctx context.Context) (*ListProductsWithBonusesResponse, error) {
	path := "/api/client/products_with_bonuses"
	var resp ListProductsWithBonusesResponse
	if err := s.Client.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
