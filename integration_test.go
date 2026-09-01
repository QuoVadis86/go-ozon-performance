package ozon

import (
	"context"
	"os"
	"testing"
)

func skipIfNoCreds(t *testing.T) {
	t.Helper()
	if os.Getenv("OZON_CLIENT_ID") == "" || os.Getenv("OZON_CLIENT_SECRET") == "" {
		t.Skip("set OZON_CLIENT_ID and OZON_CLIENT_SECRET to run integration tests")
	}
}

func testClient(t *testing.T) *Client {
	t.Helper()
	return NewClient(os.Getenv("OZON_CLIENT_ID"), os.Getenv("OZON_CLIENT_SECRET"), nil)
}

func TestIntegration_APIFlow(t *testing.T) {
	skipIfNoCreds(t)
	ctx := context.Background()
	cl := testClient(t)

	t.Run("Token", func(t *testing.T) {
		tok, err := cl.Token(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if tok == "" {
			t.Fatal("empty token")
		}
		t.Logf("token: %.20s...", tok)
	})

	t.Run("CampaignList", func(t *testing.T) {
		resp, err := cl.Campaign.ListCampaigns(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("nil response")
		}
		t.Logf("campaigns: %d", len(resp.List))
		for _, c := range resp.List {
			t.Logf("  id=%s title=%s state=%s", c.ID, c.Title, c.State)
		}
	})

	t.Run("LimitsList", func(t *testing.T) {
		resp, err := cl.Campaign.GetLimitsList(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("nil response")
		}
		t.Logf("limits: %d groups", len(resp.Limits))
	})
}
