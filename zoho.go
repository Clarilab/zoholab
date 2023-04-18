package zoholab

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	reportsUri   = "https://analyticsapi.zoho.eu"
	AccountsUri  = "https://accounts.zoho.eu/oauth/v2/token"
	validJson    = "true"
	outputFormat = "JSON"
	errorFormat  = "JSON"
	apiVersion   = "1.0"
	addRow       = "ADDROW"
)

// ZohoService is the struct for the zoho service.
type ZohoService struct {
	clientID     string
	clientSecret string
	refreshToken string
	restyClient  *resty.Client
}

// NewZohoService instantiates a new zoho service.
func NewZohoService(restyClient *resty.Client, clientID, clientSecret, refreshtoken string) *ZohoService {
	return &ZohoService{
		restyClient:  restyClient,
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshtoken,
	}
}

// GetUri returns the URI for the specified table in the Zoho Analytics account.
func (s *ZohoService) GetUri(emailID, dbName, tbName string) string {
	// Join path is not used because the path must not be escaped.
	return fmt.Sprintf("%s/api/%s/%s/%s", reportsUri, emailID, urlSplCharReplace(dbName), urlSplCharReplace(tbName))
}

// AddRow adds row to the specified table identified by the URI.
func (s *ZohoService) AddRow(tableUri string, columnValues map[string]string) (interface{}, error) {
	const errMessage = "could not add row in zoho"

	addedRows, err := s.sendAPIRequest(columnValues, true, tableUri, addRow)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}
	return addedRows, nil
}

// SendAPIRequest sends a request to the zoho api.
func (s *ZohoService) sendAPIRequest(config map[string]string, isreturn bool, path, action string) (interface{}, error) {
	const errMsg = "could not send api request"

	resp, err := s.restyClient.
		R().
		SetHeader("User-Agent", "ZohoAnalytics GoLibrary").
		SetQueryParams(map[string]string{
			"ZOHO_ACTION":        action,
			"ZOHO_OUTPUT_FORMAT": outputFormat,
			"ZOHO_ERROR_FORMAT":  errorFormat,
			"ZOHO_API_VERSION":   apiVersion,
			"ZOHO_VALID_JSON":    validJson,
		}).
		SetQueryParams(config).
		Post(path)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	if resp.IsError() {
		return nil, errors.Wrap(FillApiError(resp.Body()), errMsg)
	}

	return resp, nil
}

// Internally used. For handling special character's in the workspace name or table name.
func urlSplCharReplace(value string) string {
	value = strings.Replace(value, "/", "(/)", -1)
	return strings.Replace(value, "\\", "(//)", -1)
}
