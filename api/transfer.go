package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/WanCodeBase/GinModule/db/sqlc"
	"github.com/WanCodeBase/GinModule/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferReq struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	fromAccount, ok := server.validateCurrency(ctx, req.FromAccountId, req.Currency)
	if !ok {
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != payload.Username {
		err := errors.New("account owner is not match")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if _, ok := server.validateCurrency(ctx, req.ToAccountId, req.Currency); !ok {
		return
	}

	result, err := server.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: req.FromAccountId,
		ToAccountID:   req.ToAccountId,
		Amount:        req.Amount,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest,
			fmt.Sprintf("Currency cannot match: %s vs %s", account.Currency, currency))
		return account, false
	}
	return account, true
}
