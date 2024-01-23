package cdcexchange

import (
	"context"
	"fmt"

	"github.com/sngyai/go-cryptocom/internal/api"
	"github.com/sngyai/go-cryptocom/internal/auth"
)

const (
	methodGetDepositAddress = "private/get-deposit-address"
)

type (
	// GetDepositAddressRequest is the request params sent for the private/get-deposit-address API.
	//
	// The maximum duration between Start and EndTime is 24 hours.
	//
	// You will receive an INVALID_DATE_RANGE error if the difference exceeds the maximum duration.
	//
	// For users looking to pull longer historical deposit data, users can create a loop to make a request
	// for each 24-period from the desired start to end time.
	GetDepositAddressRequest struct {
		// Currency represents the currency symbol for the deposits (e.g. BTC or ETH).
		// if Currency is omitted, all currencies will be returned.
		Currency string `json:"currency"`
	}

	// GetDepositAddressResponse is the base response returned from the private/get-deposit-address API.
	GetDepositAddressResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result GetDepositAddressResult `json:"result"`
	}

	// GetDepositAddressResult is the result returned from the private/get-deposit-address API.
	GetDepositAddressResult struct {
		// DepositList is the array of deposits.
		DepositAddressList []DepositAddress `json:"deposit_address_list"`
	}

	DepositAddress struct {
		Currency   string `json:"currency"`
		CreateTime int64  `json:"create_time"`
		Id         string `json:"id"`
		Address    string `json:"address"`
		Status     string `json:"status"`
		Network    string `json:"network"`
	}
)

// GetDepositAddress gets the deposit address for a particular instrument.
//
// Pagination is handled using page size (Default: 20, Max: 200) & number (0-based).
// If paging is used, enumerate each page (starting with 0) until an empty deposit_list array appears in the response.
//
// req.Timeframe can be left blank to get deposits for all instruments.
//
// Method: private/get-deposit-address
func (c *Client) GetDepositAddress(ctx context.Context, req GetDepositAddressRequest) ([]DepositAddress, error) {
	var (
		id        = c.idGenerator.Generate()
		timestamp = c.clock.Now().UnixMilli()
		params    = make(map[string]interface{})
	)

	if req.Currency != "" {
		params["currency"] = req.Currency
	}

	signature, err := c.signatureGenerator.GenerateSignature(auth.SignatureRequest{
		APIKey:    c.apiKey,
		SecretKey: c.secretKey,
		ID:        id,
		Method:    methodGetDepositAddress,
		Timestamp: timestamp,
		Params:    params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	body := api.Request{
		ID:        id,
		Method:    methodGetDepositAddress,
		Nonce:     timestamp,
		Params:    params,
		Signature: signature,
		APIKey:    c.apiKey,
	}

	var GetDepositAddressResponse GetDepositAddressResponse
	statusCode, err := c.requester.Post(ctx, body, methodGetDepositAddress, &GetDepositAddressResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, GetDepositAddressResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return GetDepositAddressResponse.Result.DepositAddressList, nil
}
