package api

import (
	"database/sql"
	"fmt"
	db "github.com/WanCodeBase/GinModule/db/sqlc"
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

	if !server.validateCurrency(ctx, req.FromAccountId, req.Currency) {
		return
	}
	if !server.validateCurrency(ctx, req.ToAccountId, req.Currency) {
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

func (server *Server) validateCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	if account.Currency != currency {
		ctx.JSON(http.StatusBadRequest,
			fmt.Sprintf("Currency cannot match: %s vs %s", account.Currency, currency))
		return false
	}
	return true
}
