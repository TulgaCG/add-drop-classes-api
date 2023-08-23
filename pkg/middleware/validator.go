package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
)

func Validator(v *validator.Validate) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.ValidatorCtxKey, v)
		c.Next()
	}
}
