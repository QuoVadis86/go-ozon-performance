package product

import (
	"context"
	"strings"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

type Service struct{ Client *transport.Client }

// 将商品添加到广告活动
func (s *Service) AddProducts(ctx context.Context, campaignId transport.Uint64, req *AddProductsRequest) error {
	path := "/api/client/campaign/{campaignId}/products"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	return s.Client.Post(ctx, path, req, nil)
}

// 从活动中删除商品
func (s *Service) DeleteProducts(ctx context.Context, campaignId transport.Uint64, req *DeleteProductsRequest) error {
	path := "/api/client/campaign/{campaignId}/products/delete"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	return s.Client.Post(ctx, path, req, nil)
}

// 商品的竞争性投注
func (s *Service) GetProductsCompetitiveBids(ctx context.Context, campaignId string, req *GetProductsCompetitiveBidsRequest) (*GetProductsCompetitiveBidsResponse, error) {
	path := "/api/client/campaign/{campaignId}/products/bids/competitive"
	path = strings.Replace(path, "{campaignId}", campaignId, 1)
	var resp GetProductsCompetitiveBidsResponse
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 活动商品列表
func (s *Service) GetProductsV2(ctx context.Context, campaignId transport.Uint64, req *GetProductsV2Request) (*GetProductsResponseV2, error) {
	path := "/api/client/campaign/{campaignId}/v2/products"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	var resp GetProductsResponseV2
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 更新商品价格
func (s *Service) UpdateProducts(ctx context.Context, campaignId transport.Uint64, req *UpdateProductsRequest) error {
	path := "/api/client/campaign/{campaignId}/products"
	path = strings.Replace(path, "{campaignId}", campaignId.String(), 1)
	return s.Client.Put(ctx, path, req, nil)
}
