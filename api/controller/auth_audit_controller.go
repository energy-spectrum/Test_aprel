package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

type AuditRequest struct {
	Token string `form:"token" binding:"required"`
}

type AuditResponse struct {
	AuditEvents []db.AuditEvent `json:"auditEvents"`
}

// @Summary Authorize
// @Description Authenticate a user with login and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login formData string true "login of the user"
// @Param password formData string true "Password for the user account"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /authorize [post]
func (ac *AuthAuditController) GetAuditEvents(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse("X-Token header is missing"))
		return
	}

	expirationTime, err := ac.Store.SessionRepo.CheckToken(ctx, token)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			ctx.JSON(http.StatusUnauthorized, errorResponse("Invalid token"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	userIDStr, err := util.ExtractIDFromToken(token, expirationTime.String())
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err.Error()))
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	events, err := ac.Store.AuthAuditRepo.GetEvents(ctx, userID)
	if err != nil {
		//Think about this logic
		if errors.Is(err, util.ErrNotFound) {
			logrus.Printf("not found by user id %d", userID)
			ctx.JSON(http.StatusOK, AuditResponse{
				AuditEvents: []db.AuditEvent{},
			})
			return
		}
		logrus.Errorf("failed to get events by user id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Sprintf("failed to get events by user id %d", userID)))
		return
	}

	ctx.JSON(http.StatusOK, AuditResponse{
		AuditEvents: events,
	})
}

// @Summary Authorize
// @Description Authenticate a user with login and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login formData string true "login of the user"
// @Param password formData string true "Password for the user account"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /authorize [post]
func (ac *AuthAuditController) ClearAudit(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse("X-Token header is missing"))
		return
	}

	expirationTime, err := ac.Store.SessionRepo.CheckToken(ctx, token)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			ctx.JSON(http.StatusUnauthorized, errorResponse("Invalid token"))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	userIDStr, err := util.ExtractIDFromToken(token, expirationTime.String())
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err.Error()))
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	//Clear audit of user by user id
	err = ac.Store.AuthAuditRepo.ClearAuditByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, successResponse("Cleared audit of user"))
}
