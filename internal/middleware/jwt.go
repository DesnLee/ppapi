package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ppapi.desnlee.com/internal/jwt_helper"
	"ppapi.desnlee.com/internal/model"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authStr := c.GetHeader("Authorization")
		if authStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.MsgResponse{
				Msg: "未携带 token",
			})
			return
		}

		parts := strings.SplitN(authStr, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.MsgResponse{
				Msg: "token 格式错误",
			})
			return
		}

		claims, err := jwt_helper.ParseJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.MsgResponse{
				Msg: "token 无效",
			})
			return
		}

		if claims.UserID == uuid.Nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.MsgResponse{
				Msg: "token 无效",
			})
			return
		}
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
