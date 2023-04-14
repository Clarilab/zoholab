package middlewares

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_isAccesTokenValid(t *testing.T) {
	now := time.Now()

	oneHour := time.Hour.Seconds()
	oneMinute := time.Minute.Seconds()

	oneHourOffset := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	oneMinuteOffset := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()-1, now.Second(), now.Nanosecond(), now.Location())

	testcases := map[string]struct {
		expectedResult bool
		middleware     *AuthTokenMiddleware
	}{
		"happy path": {
			expectedResult: true,
			middleware: &AuthTokenMiddleware{
				accessToken: &AuthToken{
					AccessToken: "asdasdasd",
					ExpiresIn:   int(oneHour),
				},
				lastRequest: &oneMinuteOffset,
			},
		},
		"bad path": {
			expectedResult: false,
			middleware: &AuthTokenMiddleware{
				accessToken: &AuthToken{
					AccessToken: "asdasdasd",
					ExpiresIn:   int(oneHour),
				},
				lastRequest: &oneHourOffset,
			},
		},
		"bad path 2": {
			expectedResult: false,
			middleware: &AuthTokenMiddleware{
				accessToken: &AuthToken{
					AccessToken: "asdasdasd",
					ExpiresIn:   int(oneMinute),
				},
				lastRequest: &oneHourOffset,
			},
		},
		"bad path 3": {
			expectedResult: false,
			middleware: &AuthTokenMiddleware{
				accessToken: &AuthToken{
					AccessToken: "",
					ExpiresIn:   int(oneHour),
				},
				lastRequest: &oneMinuteOffset,
			},
		},
		"bad path 4": {
			expectedResult: false,
			middleware:     &AuthTokenMiddleware{},
		},
		"bad path 5": {
			expectedResult: false,
			middleware: &AuthTokenMiddleware{
				accessToken: &AuthToken{
					AccessToken: "asdasda",
					ExpiresIn:   int(oneHour),
				},
				lastRequest: nil,
			},
		},
		"bad path 6": {
			expectedResult: false,
			middleware: &AuthTokenMiddleware{
				accessToken: &AuthToken{
					AccessToken: "asdasda",
					ExpiresIn:   int(oneHour),
				},
				lastRequest: &time.Time{},
			},
		},
	}

	for name, tc := range testcases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			result := tc.middleware.isAccesTokenValid()

			assert.Equal(t, tc.expectedResult, result)
		})

	}
}
