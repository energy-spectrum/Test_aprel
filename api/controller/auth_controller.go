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

type AuthController struct {
	Store db.Store
	Env   *bootstrap.Env
}

type AuthRequest struct {
	Login    string `form:"login" binding:"required,min=1,max=25"`
	Password string `form:"password" binding:"required,min=1,max=25"`
}

type AuthResponse struct {
	Token string `json:"accessToken"`
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
// @Router auth/authorize [post]
func (ac *AuthController) Authorize(ctx *gin.Context) {
	var request AuthRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	userID, blocked, err := ac.Store.UserRepo.GetUserIDAndBlocked(ctx, request.Login, request.Password)
	if err != nil {
		if err == util.ErrInvalidPassword {
			//Write event Invalid Password
			err = ac.Store.AuthAuditRepo.WriteEvent(ctx, userID, db.InvalidPassword)
			if err != nil {
				logrus.Errorf("failed to write event: %v", err)
			}

			failedLoginAttemts, err := ac.Store.UserRepo.IncrementFailedLoginAttempts(ctx, userID)
			if err != nil {
				logrus.Errorf("failed to increment failed login attempts: %v", err)
				ctx.JSON(http.StatusInternalServerError, errorResponse("Invalid login or password"))
				return
			}

			if failedLoginAttemts > 4 {
				//Block user by id
				err = ac.Store.UserRepo.Block(ctx, userID)
				if err != nil {
					logrus.Errorf("failed to block user by id %d: %v", userID, err)
					ctx.JSON(http.StatusInternalServerError, errorResponse("Invalid login or password"))
					return
				}

				//Write event Block
				err = ac.Store.AuthAuditRepo.WriteEvent(ctx, userID, db.Block)
				if err != nil {
					logrus.Errorf("failed to write event: %v", err)
				}

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

	if blocked {
		ctx.JSON(http.StatusForbidden, errorResponse("User is blocked"))
		return
	}

	const lifetime = time.Hour * 24 * 7
	expirationTime := time.Now().Add(lifetime)
	expirationTime = expirationTime.Truncate(100 * time.Microsecond)
	//logrus.Printf("userID: %d, expirationTime: %s", userID, expirationTime.String())
	token,err := util.CreateAccessToken(userID, expirationTime.String(), ac.Env.TokenExpiryHour)
	if err != nil {
		logrus.Errorf("failed to create access token: %v", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to create access token")
		return
	}
	//Save access token
	err = ac.Store.SessionRepo.SaveToken(ctx, token, expirationTime)
	if err != nil {
		logrus.Errorf("failed to save token: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, "Error on the server")
		return
	}

	// Write event Login
	err = ac.Store.AuthAuditRepo.WriteEvent(ctx, userID, db.Login)
	if err != nil {
		logrus.Errorf("failed to write event: %v", err.Error())
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		Token: token,
	})
}
