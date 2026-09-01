package ozon

import (
	"context"

	"github.com/QuoVadis86/go-ozon-performance/ad"
	"github.com/QuoVadis86/go-ozon-performance/campaign"
	"github.com/QuoVadis86/go-ozon-performance/product"
	"github.com/QuoVadis86/go-ozon-performance/searchpromo"
	"github.com/QuoVadis86/go-ozon-performance/statistics"
	"github.com/QuoVadis86/go-ozon-performance/transport"
	"github.com/QuoVadis86/go-ozon-performance/vendor"
)

// ClientOptions is an alias for transport.Options.
type ClientOptions = transport.Options

// Client aggregates the Performance API service modules.
type Client struct {
	Campaign    *campaign.Service
	Statistics  *statistics.Service
	Ad          *ad.Service
	Product     *product.Service
	SearchPromo *searchpromo.Service
	Vendor      *vendor.Service

	transport *transport.Client
}

// NewClient creates a Performance API client.
//
// clientID and clientSecret are obtained from the seller personal cabinet:
// Settings → API keys → Performance API. The client fetches the access
// token automatically and refreshes it transparently once expired.
func NewClient(clientID, clientSecret string, opts *ClientOptions) *Client {
	ic := transport.New(clientID, clientSecret, opts)
	return &Client{
		Campaign:    &campaign.Service{Client: ic},
		Statistics:  &statistics.Service{Client: ic},
		Ad:          &ad.Service{Client: ic},
		Product:     &product.Service{Client: ic},
		SearchPromo: &searchpromo.Service{Client: ic},
		Vendor:      &vendor.Service{Client: ic},
		transport:   ic,
	}
}

// Token returns the current access token, fetching it if necessary.
func (c *Client) Token(ctx context.Context) (string, error) {
	return c.transport.Token(ctx)
}
