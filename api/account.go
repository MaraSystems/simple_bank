package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/MaraSystems/simple_bank/db/sqlc"
	"github.com/MaraSystems/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type IDRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authorizationPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "accounts_owner_fkey", "owner_currency_key", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func getAccountService(ctx *gin.Context, store db.Store, id int64) (db.Account, error) {
	account, err := store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, err
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, err
	}

	authorizationPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authorizationPayload.Username != account.Owner {
		err = fmt.Errorf("Account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return account, err
	}

	return account, nil
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := getAccountService(ctx, server.store, req.ID)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1,max=20"`
	Offset int32 `form:"offset" binding:"required,min=0"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authorizationPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authorizationPayload.Username,
		Limit:  req.Limit,
		Offset: req.Offset * req.Limit,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req IDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := getAccountService(ctx, server.store, req.ID)
	if err != nil {
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")
}
