package product

import (
	"context"
	"net/http"
	"testing"

	"github.com/QuoVadis86/go-ozon-performance/transport"
)

var testCtx = context.Background()

func TestAddProducts(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	err := svc.AddProducts(testCtx, transport.Uint64(0), &AddProductsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteProducts(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	err := svc.DeleteProducts(testCtx, transport.Uint64(0), &DeleteProductsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetProductsCompetitiveBids(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GetProductsCompetitiveBids(testCtx, "", &GetProductsCompetitiveBidsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetProductsV2(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	_, err := svc.GetProductsV2(testCtx, transport.Uint64(0), &GetProductsV2Request{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateProducts(t *testing.T) {
	cl, srv := transport.NewTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}, "test-token")
	defer srv.Close()
	svc := &Service{Client: cl}
	err := svc.UpdateProducts(testCtx, transport.Uint64(0), &UpdateProductsRequest{})
	if err != nil {
		t.Fatal(err)
	}
}
