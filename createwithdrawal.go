package cdcexchange

import (
	"context"
	"fmt"

	"github.com/sngyai/go-cryptocom/internal/api"
	"github.com/sngyai/go-cryptocom/internal/auth"
)

const (
	methodCreateWithdrawal = "private/create-withdrawal"
)

type (
	// CreateWithdrawalRequest is the request params sent for the private/create-withdrawal API.
	//
	// The maximum duration between Start and EndTime is 24 hours.
	//
	// You will receive an INVALID_DATE_RANGE error if the difference exceeds the maximum duration.
	//
	// For users looking to pull longer historical withdrawal data, users can create a loop to make a request
	// for each 24-period from the desired start to end time.
	CreateWithdrawalRequest struct {
		// Currency represents the currency symbol for the withdrawals (e.g. BTC or ETH).
		// if Currency is omitted, all currencies will be returned.
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
		Address  string  `json:"address"`

		ClientWid  string `json:"client_wid"`
		AddressTag string `json:"address_tag"`
		NetworkId  string `json:"network_id"`
	}

	// CreateWithdrawalResponse is the base response returned from the private/create-withdrawal API.
	CreateWithdrawalResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result CreateWithdrawalResult `json:"result"`
	}

	// CreateWithdrawalResult is the result returned from the private/create-withdrawal API.
	CreateWithdrawalResult struct {
		Id         int64   `json:"id"`
		Amount     float64 `json:"amount"`
		Fee        float64 `json:"fee"`
		Symbol     string  `json:"symbol"`
		Address    string  `json:"address"`
		ClientWid  string  `json:"client_wid"`
		CreateTime int64   `json:"create_time"`
		NetworkId  string  `json:"network_id"`
	}
)

// CreateWithdrawal gets the withdrawal history for a particular instrument.
//
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// If paging is used, enumerate each page (starting with 0) until an empty withdrawal_list array appears in the response.
//
// req.Timeframe can be left blank to get withdrawals for all instruments.
//
// Method: private/create-withdrawal
func (c *Client) CreateWithdrawal(ctx context.Context, req CreateWithdrawalRequest) (*CreateWithdrawalResult, error) {
	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.Currency != "" {
		params["currency"] = req.Currency
	}
	if req.ClientWid != "" {
		params["client_wid"] = req.ClientWid
	}
	if req.Amount != 0 {
		params["amount"] = req.Amount
	}
	if req.Address != "" {
		params["address"] = req.Address
	}

	if req.AddressTag != "" {
		params["address_tag"] = req.AddressTag
	}
	if req.NetworkId != "" {
		params["network_id"] = req.NetworkId
	}

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodCreateWithdrawal,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodCreateWithdrawal,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var CreateWithdrawalResponse CreateWithdrawalResponse
	statusCode, err := c.requester.Post(ctx, body, methodCreateWithdrawal, &CreateWithdrawalResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, CreateWithdrawalResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return &CreateWithdrawalResponse.Result, nil
}
