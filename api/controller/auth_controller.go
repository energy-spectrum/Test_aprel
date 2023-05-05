package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/santosh/gingo/docs"
	"github.com/sirupsen/logrus"

	"app/config"
	"app/db"
	"app/internal/response"
	"app/internal/util"
)

type AuthController struct {
	Store db.Store
	Env   *config.Env
}

type AuthRequest struct {
	Login    string `form:"login" binding:"required,min=2,max=25"`
	Password string `form:"password" binding:"required,min=2,max=25"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// @Summary Authorize
// @Description It authorizes a user with login and password. Returns a token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "User credentials (login and password)."
// @Success 200 {object} AuthResponse
// @Failure 400 {object} response.errorResponse
// @Failure 401 {object} response.errorResponse
// @Failure 404 {object} response.errorResponse
// @Failure 500 {object} response.errorResponse
// @Router /auth/authorize [post]
func (ac *AuthController) Authorize(ctx *gin.Context) {
	var request AuthRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	//Get user id and blocked by login and hashed password.
	hashedPassword := util.HashPassword(request.Password)
	userID, isUserblocked, err := ac.Store.UserRepo.GetUserIDAndBlocked(ctx, request.Login, hashedPassword)
	if err != nil {
		switch err {
		case util.ErrInvalidPassword:
			ac.handleInvalidPassword(ctx, userID)
		case util.ErrNotFound:
			ctx.JSON(http.StatusNotFound, response.ErrorResponse("Invalid login or password"))
		default:
			logrus.Errorf("failed to get userID and blocked: %v", err.Error())
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to get userID and blocked"))
		}
		return
	}

	if isUserblocked {
		ctx.JSON(http.StatusForbidden, response.ErrorResponse("User is blocked"))
		return
	}

	token, err := ac.generateAndSaveToken(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	ac.writeEvent(ctx, userID, db.Login)

	ctx.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}

func (ac *AuthController) handleInvalidPassword(ctx *gin.Context, userID int) {
	ac.writeEvent(ctx, userID, db.InvalidPassword)

	failedLoginAttempts, err := ac.Store.UserRepo.IncrementFailedLoginAttempts(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to increment failed login attempts: %v", err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse("Invalid login or password"))
		return
	}

	if failedLoginAttempts > ac.Env.MaxFailedLoginAttempts {
		ac.handleUserBlock(ctx, userID)
		return
	}

	ctx.JSON(http.StatusNotFound, response.ErrorResponse("Invalid login or password"))
}

func (ac *AuthController) handleUserBlock(ctx *gin.Context, userID int) {
	err := ac.Store.UserRepo.Block(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to block user by id %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse("Failed to block user"))
		return
	}

	ac.writeEvent(ctx, userID, db.Block)

	ctx.JSON(http.StatusForbidden, response.ErrorResponse("Too many failed login attempts"))
}

func (ac *AuthController) writeEvent(ctx *gin.Context, userID int, event db.EventType) {
	err := ac.Store.AuthAuditRepo.WriteEvent(ctx, userID, event)
	if err != nil {
		logrus.Errorf("failed to write event: %v", err.Error())
	}
}

func (ac *AuthController) generateAndSaveToken(ctx *gin.Context, userID int) (string, error) {
	token := util.GenerateUniqueToken()
	expirationTime := util.GetExpirationTimeForToken(ac.Env)

	err := ac.Store.SessionRepo.SaveToken(ctx, token, expirationTime, userID)
	if err != nil {
		logrus.Errorf("failed to save token: %v", err.Error())
		return "", util.ErrFailedToSaveToken
	}

	return token, nil
}
