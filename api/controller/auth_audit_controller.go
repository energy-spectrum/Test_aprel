package controller

import (
	"net/http"
	_ "time"

	"github.com/gin-gonic/gin"
	_ "github.com/santosh/gingo/docs"
	"github.com/sirupsen/logrus"

	"app/config"
	"app/db"
	"app/internal/response"
)

type AuthAuditController struct {
	Store db.Store
	Env   *config.Env
}

type GetAuthAuditResponse struct {
	AuthAudit []db.AuthAuditEvent `json:"AuthAudit"`
}

// @Summary GetAuthAudit
// @Description Get auth audit of user.
// @Tags Auth-Audit
// @Accept json
// @Produce json
// @Success 200 {object} GetAuthAuditResponse
// @Failure 400 {object} response.errorResponse
// @Failure 401 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /auth-audit/history [get]
// @security ApiKeyAuth
func (ac *AuthAuditController) GetAuthAudit(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	events, err := ac.Store.AuthAuditRepo.GetAuthAuditByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to get auth audit by user id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to get auth audit by user id"))
		return
	}

	ctx.JSON(http.StatusOK, GetAuthAuditResponse{
		AuthAudit: events,
	})
}

// @Summary ClearAuthAudit
// @Description Clear auth audit of user.
// @Tags Auth-Audit
// @Accept json
// @Produce json
// @Success 200 {object} response.successResponse
// @Failure 400 {object} response.errorResponse
// @Failure 401 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /auth-audit/clear [delete]
// @security ApiKeyAuth
func (ac *AuthAuditController) ClearAuthAudit(ctx *gin.Context) {
	userID := ctx.GetInt("userID")
	err := ac.Store.AuthAuditRepo.ClearAuthAuditByUserID(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to clear auth audit by user id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to clear auth audit of user"))
		return
	}

	err = ac.Store.UserRepo.Unblock(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to unblock by user id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to clear auth audit of user"))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse("Auth audit of user cleared"))
}
