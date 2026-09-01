package campaign

import (
	"encoding/json"
	"testing"
)

// 验证 swagger 示例响应的数字形式 id/budget 与字符串形式都能解码到 Uint64 字段。
func TestSwaggerExampleDecode(t *testing.T) {
	raw := `{
		"list": [{
			"id": 48949,
			"paymentType": "CPM",
			"title": "横幅营销活动",
			"state": "CAMPAIGN_STATE_RUNNING",
			"advObjectType": "BANNER",
			"fromDate": "2019-10-07",
			"toDate": "2021-10-07",
			"dailyBudget": 504000000,
			"placement": ["PLACEMENT_PDP"],
			"budget": 50000000,
			"createdAt": "2019-10-07T06:28:44.055042Z",
			"updatedAt": "2020-10-01T06:28:44.055042Z",
			"productAutopilotStrategy": "NO_AUTO_STRATEGY",
			"productCampaignMode": "PRODUCT_CAMPAIGN_MODE_AUTO"
		}]
	}`
	var resp CampaignsList
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.List) != 1 {
		t.Fatalf("len = %d", len(resp.List))
	}
	c := resp.List[0]
	if c.ID != 48949 {
		t.Fatalf("id = %d", c.ID)
	}
	if c.Budget != 50000000 {
		t.Fatalf("budget = %d", c.Budget)
	}
	if err := json.Unmarshal([]byte(`{"list":[{"id":"777"}]}`), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.List[0].ID != 777 {
		t.Fatalf("id = %d", resp.List[0].ID)
	}
}
