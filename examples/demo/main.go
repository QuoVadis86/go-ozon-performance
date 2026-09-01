// Command demo 演示 go-ozon-performance 的基本用法。
//
// 运行: OZON_CLIENT_ID=xxx OZON_CLIENT_SECRET=xxx go run ./examples/demo
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/QuoVadis86/go-ozon-performance"
	"github.com/QuoVadis86/go-ozon-performance/campaign"
	"github.com/QuoVadis86/go-ozon-performance/statistics"
	"github.com/QuoVadis86/go-ozon-performance/transport"
)

func main() {
	clientID := os.Getenv("OZON_CLIENT_ID")
	secret := os.Getenv("OZON_CLIENT_SECRET")
	if clientID == "" || secret == "" {
		log.Fatal("set OZON_CLIENT_ID and OZON_CLIENT_SECRET")
	}

	cl := ozon.NewClient(clientID, secret, nil)
	ctx := context.Background()

	// 1. 活动列表
	camps, err := cl.Campaign.ListCampaigns(ctx, &campaign.ListCampaignsRequest{
		AdvObjectType: "SKU",
		Page:          1,
		PageSize:      10,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("campaigns: %d\n", len(camps.List))
	for _, c := range camps.List {
		fmt.Printf("  id=%s title=%s state=%s\n", c.ID, c.Title, c.State)
	}

	// 2. 异步统计报告流程
	if len(camps.List) == 0 {
		return
	}
	sub, err := cl.Statistics.SubmitRequest(ctx, &statistics.StatisticsRequest{
		Campaigns: []transport.Uint64{camps.List[0].ID},
		DateFrom:  time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		DateTo:    time.Now().Format("2006-01-02"),
		GroupBy:   statistics.GroupByDate,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("report UUID: %s\n", sub.UUID)

	// 3. 轮询状态
	for i := 0; i < 12; i++ {
		st, err := cl.Statistics.StatisticsCheck(ctx, sub.UUID)
		if err != nil {
			log.Fatal(err)
		}
		if st.State == statistics.StatisticsRequestStateOK {
			fmt.Printf("report ready: %s\n", st.Link)
			break
		}
		fmt.Printf("state: %s, waiting...\n", st.State)
		time.Sleep(5 * time.Second)
	}
}
