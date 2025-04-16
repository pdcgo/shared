package gorm_commenter

import (
	"github.com/gin-gonic/gin"
)

func GormCommenterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		route := ctx.Request.URL.Path
		ctx.Set("route", route)

		ctx.Next()
	}
}
