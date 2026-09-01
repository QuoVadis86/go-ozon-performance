package searchpromo

import (
	"context"
	"net/http"
	"testing"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

var testCtx = context.Background()

func TestActivateAllSkuPromoCampaign(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.ActivateAllSkuPromoCampaign(testCtx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBatchDisableProducts(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	err := svc.BatchDisableProducts(testCtx, &BatchDisableProductsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBatchEnableProducts(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.BatchEnableProducts(testCtx, &BatchEnableProductsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeactivateAllSkuPromoCampaign(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	err := svc.DeactivateAllSkuPromoCampaign(testCtx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteSearchPromoBidsV2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	err := svc.DeleteSearchPromoBidsV2(testCtx, &DeleteSearchPromoBidsRequestV2{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetProductsRecommendedBids(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GetProductsRecommendedBids(testCtx, &GetProductsRecommendedBidsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListSearchPromoProductsV2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.ListSearchPromoProductsV2(testCtx, &ListSearchPromoProductsRequestV2{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetSearchPromoBidsV2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.SetSearchPromoBidsV2(testCtx, &SetSearchPromoBidsRequestV2{})
	if err != nil {
		t.Fatal(err)
	}
}
