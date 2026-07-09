package ginredoc

import (
	"github.com/gin-gonic/gin"
	"github.com/VendSYSTEM/go-redoc/v2"
)

func New(doc redoc.Redoc) gin.HandlerFunc {
	handle := doc.Handler()
	return func(ctx *gin.Context) {
		handle(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}
