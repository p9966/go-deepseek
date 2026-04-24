package deepseek

import (
	"context"
	"net/http"
)

type BalanceResponse struct {
	IsAvailable  bool          `json:"is_available"`
	BalanceInfos []BalanceInfo `json:"balance_infos"`
}

type BalanceInfo struct {
	Currency        string `json:"currency"`
	TotalBalance    string `json:"total_balance"`
	GrantedBalance  string `json:"granted_balance"`
	ToppedUpBalance string `json:"topped_up_balance"`
}

const BalanceSuffix = "/user/balance"

func (c *Client) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	request, err := c.newRequest(ctx, http.MethodGet, BalanceSuffix, nil)
	if err != nil {
		return nil, err
	}

	var result BalanceResponse
	if err := c.doJSON(request, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
