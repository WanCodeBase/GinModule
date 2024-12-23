package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WanCodeBase/GinModule/token"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	username string,
	duration time.Duration,
	authorizationType string,
) {
	authToken, err := tokenMaker.CreateToken(username, duration)
	assert.NoError(t, err)

	request.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationType, authToken))
}

func TestMiddleware(t *testing.T) {
	testCases := []struct {
		name      string
		setAuth   func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResp func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuthorization(t, request, tokenMaker, "user", time.Minute, authorizationHeaderType)
			},
			checkResp: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "ExpiredAuth",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuthorization(t, request, tokenMaker, "user", -time.Minute, authorizationHeaderType)
			},
			checkResp: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name:    "NoAuth",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {},
			checkResp: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "WrongAuthType",
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				setAuthorization(t, request, tokenMaker, "user", -time.Minute, "")
			},
			checkResp: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, nil)
				},
			)

			recoder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			assert.NoError(t, err)

			c.setAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recoder, request)
			c.checkResp(t, recoder)
		})
	}
}
