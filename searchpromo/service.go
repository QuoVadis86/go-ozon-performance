package searchpromo

import (
	"context"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

type Service struct{ Client *transport.Client }

// 启用按订单付费推广——所有商品
func (s *Service) ActivateAllSkuPromoCampaign(ctx context.Context) (*ActivateAllSkuPromoCampaignResponse, error) {
	path := "/api/client/campaign/all_sku_promo/activate"
	var resp ActivateAllSkuPromoCampaignResponse
	if err := s.Client.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 在“按订单付款”中关闭商品推广
func (s *Service) BatchDisableProducts(ctx context.Context, req *BatchDisableProductsRequest) error {
	path := "/api/client/search_promo/product/disable"
	return s.Client.Post(ctx, path, req, nil)
}

// 在“按订单付款”中开启商品推广
func (s *Service) BatchEnableProducts(ctx context.Context, req *BatchEnableProductsRequest) (*BatchEnableProductsResponse, error) {
	path := "/api/client/search_promo/product/enable"
	var resp BatchEnableProductsResponse
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 关闭按订单付费推广——所有商品
func (s *Service) DeactivateAllSkuPromoCampaign(ctx context.Context) error {
	path := "/api/client/campaign/all_sku_promo/deactivate"
	return s.Client.Get(ctx, path, nil, nil)
}

// 从“按订单付款”推广中删除商品
func (s *Service) DeleteSearchPromoBidsV2(ctx context.Context, req *DeleteSearchPromoBidsRequestV2) error {
	path := "/api/client/campaign/search_promo/v2/bids/delete"
	return s.Client.Post(ctx, path, req, nil)
}

// 推荐的商品押金
func (s *Service) GetProductsRecommendedBids(ctx context.Context, req *GetProductsRecommendedBidsRequest) (*GetProductsRecommendedBidsResponse, error) {
	path := "/api/client/search_promo/bids/recommendation"
	var resp GetProductsRecommendedBidsResponse
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// “按订单付款”推广中的商品列表
func (s *Service) ListSearchPromoProductsV2(ctx context.Context, req *ListSearchPromoProductsRequestV2) (*ListSearchPromoProductsResponseV2, error) {
	path := "/api/client/campaign/search_promo/v2/products"
	var resp ListSearchPromoProductsResponseV2
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 设置商品的价格
func (s *Service) SetSearchPromoBidsV2(ctx context.Context, req *SetSearchPromoBidsRequestV2) (*SetSearchPromoBidsResponseV2, error) {
	path := "/api/client/campaign/search_promo/v2/bids/set"
	var resp SetSearchPromoBidsResponseV2
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
