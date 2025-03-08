package deepseek

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	deepseek "github.com/p9966/go-deepseek/internal"
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
	request, err := deepseek.NewRequestBuilder(c.AuthToken).SetBaseUrl(c.BaseUrl).SetPath(BalanceSuffix).SetMethod(http.MethodGet).Build(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code:" + strconv.Itoa(resp.StatusCode))
	}

	var result BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	} else {
		return &result, nil
	}

}
