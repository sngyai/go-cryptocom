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
	methodGetDepositHistory = "private/get-deposit-history"
)

type (
	// GetDepositHistoryRequest is the request params sent for the private/get-deposit-history API.
	//
	// The maximum duration between Start and End is 24 hours.
	//
	// You will receive an INVALID_DATE_RANGE error if the difference exceeds the maximum duration.
	//
	// For users looking to pull longer historical deposit data, users can create a loop to make a request
	// for each 24-period from the desired start to end time.
	GetDepositHistoryRequest struct {
		// Currency represents the currency symbol for the deposits (e.g. BTC or ETH).
		// if Currency is omitted, all currencies will be returned.
		Currency string `json:"currency"`
		// Start is the start timestamp (milliseconds since the Unix epoch)
		// (Default: 24 hours ago)
		Start time.Time `json:"start_ts"`
		// End is the end timestamp (milliseconds since the Unix epoch)
		// (Default: now)
		End time.Time `json:"end_ts"`
		// PageSize represents maximum number of deposits returned (for pagination)
		// (Default: 20, Max: 200)
		// if PageSize is 0, it will be set as 20 by default.
		PageSize int `json:"page_size"`
		// Page represents the page number (for pagination)
		// (0-based)
		Page int `json:"page"`

		Status string `json:"status"`
	}

	// GetDepositHistoryResponse is the base response returned from the private/get-deposit-history API.
	GetDepositHistoryResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetDepositHistoryResult `json:"result"`
	}

	// GetDepositHistoryResult is the result returned from the private/get-deposit-history API.
	GetDepositHistoryResult struct {
		// DepositList is the array of deposits.
		DepositList []Deposit `json:"deposit_list"`
	}

	Deposit struct {
		Currency   string  `json:"currency"`
		Fee        float64 `json:"fee"`
		CreateTime int64   `json:"create_time"`
		Id         string  `json:"id"`
		UpdateTime int64   `json:"update_time"`
		Amount     int     `json:"amount"`
		Address    string  `json:"address"`
		Status     string  `json:"status"`
	}
)

// GetDepositHistory gets the deposit history for a particular instrument.
//
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// If paging is used, enumerate each page (starting with 0) until an empty deposit_list array appears in the response.
//
// req.InstrumentName can be left blank to get deposits for all instruments.
//
// Method: private/get-deposit-history
func (c *Client) GetDepositHistory(ctx context.Context, req GetDepositHistoryRequest) ([]Deposit, error) {
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
		Method:    methodGetDepositHistory,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetDepositHistory,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var getDepositHistoryResponse GetDepositHistoryResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetDepositHistory, &getDepositHistoryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, getDepositHistoryResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return getDepositHistoryResponse.Result.DepositList, nil
}
