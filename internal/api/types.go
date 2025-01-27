package api

import "encoding/json"

const (
	V1 = "exchange/v1/"
	V2 = "v2/"
)

type (
	Request struct {
		ID        int64                  `json:"id"`
		Method    string                 `json:"method"`
		Nonce     int64                  `json:"nonce"`
		Params    map[string]interface{} `json:"params"`
		Signature string                 `json:"sig,omitempty"`
		APIKey    string                 `json:"api_key,omitempty"`
		Version   string                 `json:"version"`
	}

	BaseResponse struct {
		ID     json.Number `json:"id"`
		Method string      `json:"method"`
		Code   json.Number `json:"code"`
	}
)
