package cdcexchange

import (
	"context"
	"fmt"
	"time"

	"github.com/sngyai/go-cryptocom/internal/api"
	"github.com/sngyai/go-cryptocom/internal/auth"
)

const (
	methodUserBalanceHistory = "private/user-balance-history"
)

type (
	UserBalance struct {
		T int64  `json:"t"`
		C string `json:"c"`
	}
	// UserBalanceHistoryRequest is the request params sent for the private/user-balance-history API.
	UserBalanceHistoryRequest struct {
		Timeframe string    `json:"timeframe"`
		EndTime   time.Time `json:"end_time"`
		Limit     int       `json:"limit"`
	}

	// UserBalanceHistoryResponse is the base response returned from the private/user-balance-history API.
	UserBalanceHistoryResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result UserBalanceHistoryResult `json:"result"`
	}

	// UserBalanceHistoryResult is the result returned from the private/user-balance-history API.
	UserBalanceHistoryResult struct {
		InstrumentName string        `json:"instrument_name"`
		Data           []UserBalance `json:"data"`
	}
)

// UserBalanceHistory gets all executed trades for a particular instrument.
// Method: private/user-balance-history
func (c *Client) UserBalanceHistory(ctx context.Context, req UserBalanceHistoryRequest) (*UserBalanceHistoryResult, error) {
	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.Timeframe != "" {
		params["timeframe"] = req.Timeframe
	}
	if req.Limit != 0 {
		params["limit"] = req.Limit
	}
	if !req.EndTime.IsZero() {
		params["end_time"] = req.EndTime.UnixMilli()
	}

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodUserBalanceHistory,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodUserBalanceHistory,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
		Version:   api.V1,
	}

	var userBalanceHistoryResponse UserBalanceHistoryResponse
	statusCode, err := c.requester.Post(ctx, body, methodUserBalanceHistory, &userBalanceHistoryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, userBalanceHistoryResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &userBalanceHistoryResponse.Result, nil
}
