package domain

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/pkg/errors"
)

// drop-in replacement for encoding/json
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ZohoError is the error that the zoho api returns.
type ZohoError struct {
	Response *Response `json:"response"`
}

// Response is part of the ZohoError response Model.
type Response struct {
	URI           string        `json:"uri"`
	Action        string        `json:"action"`
	ResponseError ResponseError `json:"error"`
}

// ResponseError is part of the ZohoError response Model.
type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Error stringifies the ServerException Error.
func (apiError ZohoError) Error() string {
	return fmt.Sprintf("Uri: %s, Action: %s, ErrorCode: %d, ErrorMessage: %s",
		apiError.Response.URI,
		apiError.Response.Action,
		apiError.Response.ResponseError.Code,
		apiError.Response.ResponseError.Message)
}

// FillApiError helper function used to fill the ZohoError model.
func FillApiError(respBody []byte) error {
	const errMsg = "could not fill api Error"

	var errResp ZohoError

	err := json.Unmarshal(respBody, &errResp)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return errResp
}
