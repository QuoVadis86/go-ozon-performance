package ad

import (
	"context"
	"net/http"
	"testing"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

var testCtx = context.Background()

func TestActivateCampaign(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.ActivateCampaign(testCtx, transport.Uint64(0), &Empty{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCalculateDynamicBudget(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.CalculateDynamicBudget(testCtx, &CalculateDynamicBudgetRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateProductCampaignCPCV2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.CreateProductCampaignCPCV2(testCtx, &CreateProductCampaignRequestV2CPC{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeactivateCampaign(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.DeactivateCampaign(testCtx, transport.Uint64(0), &Empty{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPatchProductCampaign(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.PatchProductCampaign(testCtx, transport.Uint64(0), &PatchProductCampaignRequest{})
	if err != nil {
		t.Fatal(err)
	}
}
