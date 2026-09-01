package ad

import (
	"context"
	"strings"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

type Service struct{ Client *transport.Client }

// 激活广告活动
func (s *Service) ActivateCampaign(ctx context.Context, campaignId transport.Uint64, req *Empty) (*Campaign, error) {
	path := "/api/client/campaign/{campaignId}/activate"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	var resp Campaign
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 计算最低广告活动预算
func (s *Service) CalculateDynamicBudget(ctx context.Context, req *CalculateDynamicBudgetRequest) (*CalculateDynamicBudgetResponse, error) {
	path := "/external/api/dynamic_budget"
	var resp CalculateDynamicBudgetResponse
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 创建点击付费广告活动
func (s *Service) CreateProductCampaignCPCV2(ctx context.Context, req *CreateProductCampaignRequestV2CPC) (*CampaignID, error) {
	path := "/api/client/campaign/cpc/v2/product"
	var resp CampaignID
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 关闭广告活动
func (s *Service) DeactivateCampaign(ctx context.Context, campaignId transport.Uint64, req *Empty) (*Campaign, error) {
	path := "/api/client/campaign/{campaignId}/deactivate"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	var resp Campaign
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 广告活动参数
func (s *Service) PatchProductCampaign(ctx context.Context, campaignId transport.Uint64, req *PatchProductCampaignRequest) (*CampaignID, error) {
	path := "/api/client/campaign/{campaignId}"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	var resp CampaignID
	if err := s.Client.Patch(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
