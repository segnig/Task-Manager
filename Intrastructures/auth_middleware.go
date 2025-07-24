package Intrastructures

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/segnig/task-manager/Domains"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")

		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "No Authentication header provided"})
			ctx.Abort()
			return
		}
		claims, err := ValidateToken(clientToken)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("username", claims.Username)
		ctx.Set("first_name", claims.FirstName)
		ctx.Set("last_name", claims.LastName)
		ctx.Set("user_id", claims.Uid)
		ctx.Set("user_type", claims.UserType)
		ctx.Next()
	}
}
