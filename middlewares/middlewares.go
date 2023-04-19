package middlewares

import (
	"time"

	"github.com/Clarilab/zoholab"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	accountsUri = "https://accounts.zoho.eu/oauth/v2/token"
	timeout     = 3300
)

// ZohoService is the struct for the zoho service.
type AuthTokenMiddleware struct {
	clientID     string
	clientSecret string
	accessToken  *AuthToken
	refreshToken string
	lastRequest  *time.Time
}

// NewZohoService instantiates a new zoho service.
func NewAuthTokenMiddleware(clientID, clientSecret, refreshtoken string) *AuthTokenMiddleware {
	return &AuthTokenMiddleware{
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshtoken,
	}
}

// AddAuthToken adds an oauth token to a resty request.
func (a *AuthTokenMiddleware) AddAuthTokenToRequest(client *resty.Client, request *resty.Request) error {
	const errMsg = "could not add auth token"

	token, err := a.getOAuthToken()
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	request.SetHeader("Authorization", "Zoho-oauthtoken "+token.AccessToken)

	return nil
}

// getOAuthToken gets a oauthtoken from the zoho api.
func (a *AuthTokenMiddleware) getOAuthToken() (*AuthToken, error) {
	const errMsg = "could not get oauth token"

	if a.isAccesTokenValid() {
		return a.accessToken, nil
	}

	var result AuthToken

	now := time.Now()

	a.lastRequest = &now

	resp, err := resty.New().R().
		SetQueryParams(map[string]string{
			"client_id":     a.clientID,
			"client_secret": a.clientSecret,
			"refresh_token": a.refreshToken,
			"grant_type":    "refresh_token"},
		).
		SetResult(&result).
		Post(accountsUri)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	if resp.IsError() {
		return nil, errors.Wrap(zoholab.FillApiError(resp.Body()), errMsg)
	}

	return &result, nil
}

// isAccesTokenValid checks if an acces token is valid.
func (a *AuthTokenMiddleware) isAccesTokenValid() bool {
	if a.accessToken == nil || a.accessToken.AccessToken == "" {
		return false
	}

	if a.lastRequest == nil || a.lastRequest.IsZero() {
		return false
	}

	if time.Since(*a.lastRequest).Seconds() > float64(a.accessToken.ExpiresIn-timeout) {
		return false
	}

	return true
}

// AuthToken is the model for the auth token.
type AuthToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ApiDomain   string `json:"api_domain"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
