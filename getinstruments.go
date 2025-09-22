package cdcexchange

import (
	"context"
	"fmt"

	"github.com/sngyai/go-cryptocom/internal/api"
)

const (
	methodGetInstruments = "public/get-instruments"
)

type (
	// InstrumentsResponse is the base response returned from the public/get-instruments API.
	InstrumentsResponse struct {
		// api.BaseResponse is the common response fields.
		api.BaseResponse
		// Result is the response attributes of the endpoint.
		Result InstrumentResult `json:"result"`
	}

	// InstrumentResult is the result returned from the public/get-instruments API.
	InstrumentResult struct {
		// Instruments is a list of the returned instruments.
		Instruments []Instrument `json:"data"`
	}

	// Instrument represents details of a specific currency pair
	Instrument struct {
		Symbol            string `json:"symbol"`
		InstType          string `json:"inst_type"`
		DisplayName       string `json:"display_name"`
		BaseCcy           string `json:"base_ccy"`
		QuoteCcy          string `json:"quote_ccy"`
		QuoteDecimals     int    `json:"quote_decimals"`
		QuantityDecimals  int    `json:"quantity_decimals"`
		PriceTickSize     string `json:"price_tick_size"`
		QtyTickSize       string `json:"qty_tick_size"`
		MaxLeverage       string `json:"max_leverage"`
		Tradable          bool   `json:"tradable"`
		ExpiryTimestampMs int    `json:"expiry_timestamp_ms"`
		BetaProduct       bool   `json:"beta_product"`
		UnderlyingSymbol  string `json:"underlying_symbol"`
		ContractSize      string `json:"contract_size"`
		MarginBuyEnabled  bool   `json:"margin_buy_enabled"`
		MarginSellEnabled bool   `json:"margin_sell_enabled"`
	}
)

// GetInstruments provides information on all supported instruments (e.g. BTC_USDT).
//
// Method: public/get-instruments
func (c *Client) GetInstruments(ctx context.Context) ([]Instrument, error) {
	body := api.Request{
		ID:     c.idGenerator.Generate(),
		Method: methodGetInstruments,
		Nonce:  c.clock.Now().UnixMilli(),
	}

	var instrumentsResponse InstrumentsResponse
	statusCode, err := c.requester.Get(ctx, body, methodGetInstruments, &instrumentsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to execute post request: %w", err)
	}

	if err := c.requester.CheckErrorResponse(statusCode, instrumentsResponse.Code); err != nil {
		return nil, fmt.Errorf("error received in response: %w", err)
	}

	return instrumentsResponse.Result.Instruments, nil
}
