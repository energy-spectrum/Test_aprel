package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/santosh/gingo/docs"
	"github.com/sirupsen/logrus"

	"app/bootstrap"
	"app/db"
	"app/internal/util"
)

// AuthController is responsible for handling HTTP requests and responses related to user authorization.
type AuthController struct {
	Store db.Store
	Env   *bootstrap.Env
}

// AuthRequest represents the JSON request body for the Authorize method.
type AuthRequest struct {
	Login    string `form:"login" binding:"required,min=2,max=25"`
	Password string `form:"password" binding:"required,min=2,max=25"`
}

// AuthResponse represents the JSON response body for the Authorize method.
type AuthResponse struct {
	Token string `json:"token"`
}

// Authorize handles POST requests to the /auth/authorize endpoint.
// It authorizes a user with login and password and returns a token.
// @Summary Authorize
// @Description Authorize a user with login and password. Return token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "User credentials (login and password)"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/authorize [post]
func (ac *AuthController) Authorize(ctx *gin.Context) {
	var request AuthRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	hashedPassword := util.HashPassword(request.Password)
	userID, isUserblocked, err := ac.Store.UserRepo.GetUserIDAndBlocked(ctx, request.Login, hashedPassword)
	if err != nil {
		if err == util.ErrInvalidPassword {
			//Write event Invalid Password
			ac.writeEvent(ctx, userID, db.InvalidPassword)

			//Increment failed login attempts
			failedLoginAttempts, err := ac.Store.UserRepo.IncrementFailedLoginAttempts(ctx, userID)
			if err != nil {
				logrus.Errorf("failed to increment failed login attempts: %v", err)
				ctx.JSON(http.StatusInternalServerError, errorResponse("Invalid login or password"))
				return
			}

			if failedLoginAttempts > ac.Env.MaxFailedLoginAttempts {
				//Block user by id
				err = ac.Store.UserRepo.Block(ctx, userID)
				if err != nil {
					logrus.Errorf("failed to block user by id %d: %v", userID, err)
					ctx.JSON(http.StatusInternalServerError, errorResponse("Failed to block user"))
					return
				}

				//Write event Block
				ac.writeEvent(ctx, userID, db.Block)

				ctx.JSON(http.StatusForbidden, errorResponse("Too many failed login attempts"))
				return
			}

			ctx.JSON(http.StatusNotFound, errorResponse("Invalid login or password"))
			return
		} else if err == util.ErrNotFound {
			ctx.JSON(http.StatusNotFound, errorResponse("Invalid login or password"))
			return
		}

		logrus.Errorf("failed to get userID And blocked: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, errorResponse("Error on the server"))
		return
	}

	if isUserblocked {
		ctx.JSON(http.StatusForbidden, errorResponse("User is blocked"))
		return
	}

	token, err := ac.generateAndSaveToken(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	// Write event Login
	ac.writeEvent(ctx, userID, db.Login)

	ctx.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}

func (ac *AuthController) writeEvent(ctx *gin.Context, userID int, event db.EventType) {
	err := ac.Store.AuthAuditRepo.WriteEvent(ctx, userID, event)
	if err != nil {
		logrus.Errorf("failed to write event: %v", err.Error())
	}
}

// func (ac *AuthController) blockUser(ctx *gin.Context, userID int) error {
// 	err := ac.Store.UserRepo.Block(ctx, userID)
// 	if err != nil {
// 	  logrus.Errorf("failed to block user by id %d: %v", userID, err)
// 	  return err
// 	}
// }

func (ac *AuthController) generateAndSaveToken(ctx *gin.Context, userID int) (string, error) {
	token := util.GenerateUniqueToken()
	lifetime := time.Hour * time.Duration(ac.Env.TokenExpiryHour)
	expirationTime := time.Now().Add(lifetime)
	expirationTime = expirationTime.Truncate(100 * time.Microsecond)

	err := ac.Store.SessionRepo.SaveToken(ctx, token, expirationTime, userID)
	if err != nil {
		logrus.Errorf("failed to save token: %v", err.Error())
		return "", util.ErrFailedToSaveToken
	}

	return token, nil
}
