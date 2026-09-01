package vendor

import (
	"context"
	"strings"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

type Service struct{ Client *transport.Client }

// 外部广告活动的组织标记
func (s *Service) GetVendorTag(ctx context.Context, req *GetVendorTagRequest) (*GetVendorTagResponse, error) {
	path := "/api/client/organisation/vendor_tag"
	var resp GetVendorTagResponse
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 关于UUID报告的信息
func (s *Service) VendorStatisticsCheck(ctx context.Context, uUID string, req *VendorStatisticsCheckRequest) (*VendorStatisticsResponse, error) {
	path := "/api/client/vendors/statistics/{UUID}"
	path = strings.Replace(path, "{UUID}", uUID, 1)
	var resp VendorStatisticsResponse
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 请求的报告列表，包含外部流量分析
func (s *Service) VendorStatisticsListReports(ctx context.Context, req *VendorStatisticsListReportsRequest) (*VendorStatisticsReportsList, error) {
	path := "/api/client/vendors/statistics/list"
	var resp VendorStatisticsReportsList
	if err := s.Client.Get(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// 外部流量分析报告
func (s *Service) VendorStatisticsSubmitRequest(ctx context.Context, req *VendorStatisticsRequest) (*StatisticsRequestID, error) {
	path := "/api/client/vendors/statistics"
	var resp StatisticsRequestID
	if err := s.Client.Post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
