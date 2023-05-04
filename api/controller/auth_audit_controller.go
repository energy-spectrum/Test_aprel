package controller

import (
	"errors"
	"net/http"
	_ "time"

	"github.com/gin-gonic/gin"
	_ "github.com/santosh/gingo/docs"
	"github.com/sirupsen/logrus"

	"app/bootstrap"
	"app/db"
	"app/internal/util"
)

type AuthAuditController struct {
	Store db.Store
	Env   *bootstrap.Env
}

type GetAuthAuditResponse struct {
	AuthAudit []db.AuthAuditEvent `json:"AuthAudit"`
}

// @Summary GetAuthAudit
// @Description Get auth audit of user
// @Tags Auth-Audit
// @Accept json
// @Produce json
// @Success 200 {object} GetAuthAuditResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth-audit/history [get]
// @security ApiKeyAuth
func (ac *AuthAuditController) GetAuthAudit(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse("X-Token header is missing"))
		return
	}

	userID, err := ac.Store.SessionRepo.GetUserID(ctx, token)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			ctx.JSON(http.StatusUnauthorized, errorResponse("Invalid token"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}
	//Get auth audit of user
	events, err := ac.Store.AuthAuditRepo.GetAuthAuditByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			logrus.Printf("auth audit was not found by user id %d", userID)
			ctx.JSON(http.StatusOK, GetAuthAuditResponse{
				AuthAudit: []db.AuthAuditEvent{},
			})
			return
		}
		logrus.Errorf("failed to get auth audit by user id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed toget auth audit by user id"))
		return
	}

	ctx.JSON(http.StatusOK, GetAuthAuditResponse{
		AuthAudit: events,
	})
}

// @Summary ClearAuthAudit
// @Description Clear auth audit of user
// @Tags Auth-Audit
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth-audit/clear [delete]
// @security ApiKeyAuth
func (ac *AuthAuditController) ClearAuthAudit(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse("X-Token header is missing"))
		return
	}

	userID, err := ac.Store.SessionRepo.GetUserID(ctx, token)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			ctx.JSON(http.StatusUnauthorized, errorResponse("Invalid token"))
			return
		}

		logrus.Errorf("failed to get user id by token %s: %v", token, err)
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get user id by token"))
		return
	}
	//Clear auth audit by user id
	err = ac.Store.AuthAuditRepo.ClearAuthAuditByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to clear auth audit by user id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to clear auth audit of user"))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("Auth audit of user cleared"))
}
