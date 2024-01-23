package cdcexchange

import (
	"net/http"

	"github.com/jonboulle/clockwork"

	"github.com/sngyai/go-cryptocom/errors"
	"github.com/sngyai/go-cryptocom/internal/auth"
	"github.com/sngyai/go-cryptocom/internal/id"
)

const (
	UATSandboxBaseURL = uatSandboxBaseURL
	ProductionBaseURL = productionBaseURL

	// Common API
	MethodGetInstruments = methodGetInstruments
	MethodGetBook        = methodGetBook
	MethodGetTicker      = methodGetTicker

	// Spot Trading API
	MethodGetAccountSummary = methodGetAccountSummary
	MethodCreateOrder       = methodCreateOrder
	MethodCancelOrder       = methodCancelOrder
	MethodCancelAllOrders   = methodCancelAllOrders
	MethodGetOrderHistory   = methodGetOrderHistory
	MethodGetOpenOrders     = methodGetOpenOrders
	MethodGetOrderDetail    = methodGetOrderDetail
	MethodGetTrades         = methodGetTrades
)

func (c *Client) BaseURL() string {
	return c.requester.BaseURL
}

func (c *Client) APIKey() string {
	return c.apiKey
}

func (c *Client) SecretKey() string {
	return c.secretKey
}

func (c *Client) HTTPClient() *http.Client {
	return c.requester.Client
}

func WithIDGenerator(idGenerator id.IDGenerator) ClientOption {
	return func(c *Client) error {
		if idGenerator == nil {
			return errors.InvalidParameterError{Parameter: "idGenerator", Reason: "cannot be empty"}
		}

		c.idGenerator = idGenerator
		return nil
	}
}

func WithSignatureGenerator(signatureGenerator auth.SignatureGenerator) ClientOption {
	return func(c *Client) error {
		if signatureGenerator == nil {
			return errors.InvalidParameterError{Parameter: "signatureGenerator", Reason: "cannot be empty"}
		}

		c.signatureGenerator = signatureGenerator
		return nil
	}
}

func WithClock(clock clockwork.Clock) ClientOption {
	return func(c *Client) error {
		if clock == nil {
			return errors.InvalidParameterError{Parameter: "clock", Reason: "cannot be empty"}
		}

		c.clock = clock
		return nil
	}
}

func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		if url == "" {
			return errors.InvalidParameterError{Parameter: "url", Reason: "cannot be empty"}
		}

		c.requester.BaseURL = url
		return nil
	}
}
