package zoholab

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/Clarilab/zoholab/domain"
	"github.com/Clarilab/zoholab/middlewares"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	reportsUri   = "https://analyticsapi.zoho.eu"
	validJson    = "true"
	outputFormat = "JSON"
	errorFormat  = "JSON"
	apiVersion   = "1.0"
	csvFileType  = "CSV"
	autoIdentify = "true"
	addRowAction = "ADDROW"
	importAction = "IMPORT"
)

// ZohoService is the struct for the zoho service.
type ZohoService struct {
	restyClient *resty.Client
}

// NewZohoService instantiates a new zoho service.
func NewZohoService() *ZohoService {
	return &ZohoService{
		restyClient: resty.New(),
	}
}

// SetServiceParams sets the params needed to call the zoho api.
func (s *ZohoService) SetServiceParams(clientID, clientSecret, refreshToken string) {
	authTokenMiddleware := middlewares.NewAuthTokenMiddleware(clientID, clientSecret, refreshToken)

	s.restyClient.OnBeforeRequest(authTokenMiddleware.AddAuthTokenToRequest)
}

// GetUri returns the URI for the specified table in the Zoho Analytics account.
func (s *ZohoService) GetUri(emailID, dbName, tbName string) string {
	// Join path is not used because the path must not be escaped.
	return fmt.Sprintf("%s/api/%s/%s/%s", reportsUri, emailID, urlSplCharReplace(dbName), urlSplCharReplace(tbName))
}

// AddRow adds row to the specified table identified by the URI.
func (s *ZohoService) AddRow(tableUri string, columnValues map[string]string) (*ZohoAddRowResponse, error) {
	const errMessage = "could not add row in zoho"

	var addedRows ZohoAddRowResponse

	err := s.sendAPIRequest(columnValues, nil, tableUri, addRowAction, nil, &addedRows)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	return &addedRows, nil
}

// ImportCSV import a bulk of rows as CSV
func (s *ZohoService) ImportCSV(tableUri, csvData string, config map[string]string) (*ZohoAddRowResponse, error) {
	const errMessage = "could not import csv data in zoho"

	config["ZOHO_AUTO_IDENTIFY"] = autoIdentify

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)
	_, err := bodyWriter.CreateFormFile("ZOHO_FILE", "file.csv")
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	fileReader := strings.NewReader(csvData)

	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	requestReader := io.MultiReader(bodyBuf, fileReader, closeBuf)

	headers := map[string]string{
		"Content-Type": "multipart/form-data; boundary=" + boundary,
	}

	var resp ZohoAddRowResponse
	err = s.sendAPIRequest(config, headers, tableUri, importAction, requestReader, &resp)
	if err != nil {
		return nil, errors.Wrap(err, errMessage)
	}

	return &resp, err
}

// sendAPIRequest sends a request to the zoho api.
func (s *ZohoService) sendAPIRequest(config map[string]string, additionalHeaders map[string]string, path, action string, body any, result any) error {
	const errMsg = "could not send api request"

	request := s.restyClient.
		R().
		SetQueryParams(map[string]string{
			"ZOHO_ACTION":        action,
			"ZOHO_OUTPUT_FORMAT": outputFormat,
			"ZOHO_ERROR_FORMAT":  errorFormat,
			"ZOHO_API_VERSION":   apiVersion,
			"ZOHO_VALID_JSON":    validJson,
		}).
		SetQueryParams(config).
		SetResult(result)

	if body != nil {
		request.SetBody(body)
	}

	if additionalHeaders != nil {
		request.SetHeaders(additionalHeaders)
	}

	resp, err := request.Post(path)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	if resp.IsError() {
		return errors.Wrap(domain.FillApiError(resp.Body()), errMsg)
	}

	return nil
}

// Internally used. For handling special character's in the workspace name or table name.
func urlSplCharReplace(value string) string {
	value = strings.Replace(value, "/", "(/)", -1)
	return strings.Replace(value, "\\", "(//)", -1)
}

// ZohoAddRowResponse is the response that the zoho api returns.
type ZohoAddRowResponse struct {
	Response *ZohoResponse `json:"response"`
}

// ZohoResponse is part of the ZohoAddRowResponse response Model.
type ZohoResponse struct {
	URI            string         `json:"uri"`
	Action         string         `json:"action"`
	ResponseResult ResponseResult `json:"result"`
}

// ResponseResult is part of the ZohoAddRowResponse response Model.
type ResponseResult struct {
	ColumnOrder []string   `json:"column_order"`
	Rows        [][]string `json:"rows"`
}
