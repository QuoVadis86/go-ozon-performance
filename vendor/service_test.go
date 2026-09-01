package vendor

import (
	"context"
	"net/http"
	"testing"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

var testCtx = context.Background()

func TestGetVendorTag(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GetVendorTag(testCtx, &GetVendorTagRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVendorStatisticsCheck(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.VendorStatisticsCheck(testCtx, "", &VendorStatisticsCheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVendorStatisticsListReports(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.VendorStatisticsListReports(testCtx, &VendorStatisticsListReportsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestVendorStatisticsSubmitRequest(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.VendorStatisticsSubmitRequest(testCtx, &VendorStatisticsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}
