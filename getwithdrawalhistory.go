package cdcexchange

import (
	"context"
	"fmt"
	"time"

	"github.com/sngyai/go-cryptocom/errors"
	"github.com/sngyai/go-cryptocom/internal/api"
	"github.com/sngyai/go-cryptocom/internal/auth"
)

const (
	methodGetWithdrawalHistory = "private/get-withdrawal-history"
)

type (
	// GetWithdrawalHistoryRequest is the request params sent for the private/get-withdrawal-history API.
	//
	// The maximum duration between Start and End is 24 hours.
	//
	// You will receive an INVALID_DATE_RANGE error if the difference exceeds the maximum duration.
	//
	// For users looking to pull longer historical withdrawal data, users can create a loop to make a request
	// for each 24-period from the desired start to end time.
	GetWithdrawalHistoryRequest struct {
		// Currency represents the currency symbol for the withdrawals (e.g. BTC or ETH).
		// if Currency is omitted, all currencies will be returned.
		Currency string `json:"currency"`
		// Start is the start timestamp (milliseconds since the Unix epoch)
		// (Default: 24 hours ago)
		Start time.Time `json:"start_ts"`
		// End is the end timestamp (milliseconds since the Unix epoch)
		// (Default: now)
		End time.Time `json:"end_ts"`
		// PageSize represents maximum number of withdrawals returned (for pagination)
		// (Default: 20, Max: 200)
		// if PageSize is 0, it will be set as 20 by default.
		PageSize int `json:"page_size"`
		// Page represents the page number (for pagination)
		// (0-based)
		Page int `json:"page"`

		Status string `json:"status"`
	}

	// GetWithdrawalHistoryResponse is the base response returned from the private/get-withdrawal-history API.
	GetWithdrawalHistoryResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetWithdrawalHistoryResult `json:"result"`
	}

	// GetWithdrawalHistoryResult is the result returned from the private/get-withdrawal-history API.
	GetWithdrawalHistoryResult struct {
		// WithdrawalList is the array of withdrawals.
		WithdrawalList []Withdrawal `json:"withdrawal_list"`
	}

	Withdrawal struct {
		Currency   string      `json:"currency"`
		ClientWid  string      `json:"client_wid"`
		Fee        float64     `json:"fee"`
		CreateTime int64       `json:"create_time"`
		Id         string      `json:"id"`
		UpdateTime int64       `json:"update_time"`
		Amount     float64     `json:"amount"`
		Address    string      `json:"address"`
		Status     string      `json:"status"`
		Txid       string      `json:"txid"`
		NetworkId  interface{} `json:"network_id"`
	}
)

// GetWithdrawalHistory gets the withdrawal history for a particular instrument.
//
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// If paging is used, enumerate each page (starting with 0) until an empty withdrawal_list array appears in the response.
//
// req.InstrumentName can be left blank to get withdrawals for all instruments.
//
// Method: private/get-withdrawal-history
func (c *Client) GetWithdrawalHistory(ctx context.Context, req GetWithdrawalHistoryRequest) ([]Withdrawal, error) {
	if req.PageSize < 0 {
		return nil, errors.InvalidParameterError{Parameter: "req.PageSize", Reason: "cannot be less than 0"}
	}
	if req.PageSize > 200 {
		return nil, errors.InvalidParameterError{Parameter: "req.PageSize", Reason: "cannot be greater than 200"}
	}

	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.Currency != "" {
		params["currency"] = req.Currency
	}
	if req.PageSize != 0 {
		params["page_size"] = req.PageSize
	}
	if !req.Start.IsZero() {
		params["start_ts"] = req.Start.UnixMilli()
	}
	if !req.End.IsZero() {
		params["end_ts"] = req.End.UnixMilli()
	}
	params["page"] = req.Page
	if req.Status != "" {
		params["status"] = req.Status
	}

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodGetWithdrawalHistory,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetWithdrawalHistory,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var getWithdrawalHistoryResponse GetWithdrawalHistoryResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetWithdrawalHistory, &getWithdrawalHistoryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, getWithdrawalHistoryResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return getWithdrawalHistoryResponse.Result.WithdrawalList, nil
}
