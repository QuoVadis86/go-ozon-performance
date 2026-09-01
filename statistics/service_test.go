package statistics

import (
	"context"
	"net/http"
	"testing"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

var testCtx = context.Background()

func TestAttributionSubmitRequest(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.AttributionSubmitRequest(testCtx, &StatisticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDownloadStatistics(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.DownloadStatistics(testCtx, &DownloadStatisticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateAllSkuPromoOrdersReport(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GenerateAllSkuPromoOrdersReport(testCtx, &GenerateAllSkuPromoOrdersReportRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateAllSkuPromoProductsReport(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GenerateAllSkuPromoProductsReport(testCtx, &GenerateAllSkuPromoProductsReportRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCampaignDailyStats(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GetCampaignDailyStats(testCtx, &GetCampaignDailyStatsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCampaignExpense(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GetCampaignExpense(testCtx, &GetCampaignExpenseRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListReports(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.ListReports(testCtx, &ListReportsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListReportsExternal(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.ListReportsExternal(testCtx, &ListReportsExternalRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMediaCampaignList(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.MediaCampaignList(testCtx, &MediaCampaignListRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestProductCampaignList(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.ProductCampaignList(testCtx, &ProductCampaignListRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchPromoOrdersReportSubmitRequest(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.SearchPromoOrdersReportSubmitRequest(testCtx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchPromoProductsReportSubmitRequest(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.SearchPromoProductsReportSubmitRequest(testCtx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchPromoProductsSKUStatistics2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.SearchPromoProductsSKUStatistics2(testCtx, &SearchPromoProductsSKUStatisticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatisticsCheck(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.StatisticsCheck(testCtx, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatisticsPhrasesSubmitRequest2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.StatisticsPhrasesSubmitRequest2(testCtx, &ExtstatisticsStatisticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSubmitRequest(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.SubmitRequest(testCtx, &StatisticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVideoCampaignsSubmitRequest(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.VideoCampaignsSubmitRequest(testCtx, &StatisticsVideobannerRequest{})
	if err != nil {
		t.Fatal(err)
	}
}
