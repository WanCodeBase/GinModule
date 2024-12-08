package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/WanCodeBase/GinModule/db/mock"
	db "github.com/WanCodeBase/GinModule/db/sqlc"
	"github.com/WanCodeBase/GinModule/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name      string
		accountID int64
		stubs     func(store *mockdb.MockStore)
		checkResp func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrTxDone)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: -1,
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), -1).Times(0).Return(db.Account{}, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			store := mockdb.NewMockStore(ctrl)
			c.stubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", c.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			c.checkResp(t, recorder)
		})
	}
}

func TestCreateAccountApi(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name      string
		body      gin.H
		stubs     func(store *mockdb.MockStore)
		checkResp func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
					}).
					Times(1).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "BadRequest",
			body: gin.H{
				"owner":    account.Owner,
				"currency": "AUS",
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: "AUS",
					}).
					Times(0).
					Return(account, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner":    account.Owner,
				"currency": account.Currency,
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), db.CreateAccountParams{
						Owner:    account.Owner,
						Currency: account.Currency,
					}).
					Times(1).
					Return(db.Account{}, sql.ErrTxDone)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			c.stubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/account")
			body, err := json.Marshal(c.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			c.checkResp(t, recorder)
		})
	}
}

func TestListAccountsApi(t *testing.T) {
	var n = 5
	accounts := make([]db.Account, n, n)

	for i := 0; i < n; i++ {
		accounts[i] = randomAccount()
	}

	testCases := []struct {
		name      string
		query     listAccountReq
		stubs     func(store *mockdb.MockStore)
		checkResp func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: listAccountReq{
				PageID:   1,
				PageSize: int32(n),
			},
			stubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), db.ListAccountsParams{
						Limit:  int32(n),
						Offset: 1,
					}).
					Times(1).
					Return(accounts, nil)
			},
			checkResp: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, accounts)
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			c.stubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprint(c.query.PageID))
			q.Add("page_size", fmt.Sprint(c.query.PageSize))
			request.URL.RawQuery = q.Encode()
			server.router.ServeHTTP(recorder, request)

			fmt.Println(c, recorder)
			c.checkResp(t, recorder)

		})

	}
}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount []db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccount)
}
