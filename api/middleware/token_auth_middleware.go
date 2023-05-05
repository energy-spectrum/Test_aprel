package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/db"
	"app/internal/response"
	"app/internal/util"
)

func TokenAuthMiddleware(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("X-Token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse("X-Token header is missing"))
			ctx.Abort()
			return
		}

		userID, err := store.SessionRepo.GetUserID(ctx, token)
		if err != nil {
			if errors.Is(err, util.ErrNotFound) {
				ctx.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid token"))
				ctx.Abort()
				return
			}

			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
