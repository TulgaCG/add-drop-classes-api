package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

// AuthzWithRoles middleware checks whether the user has one of the allowedRoles
// If the user has one of the roles in allowedRoles, the middleware calls the next handler
// If the user has NOT one of the roles, it returns http.StatusUnauthorized.
func AuthzWithRoles(db *database.DB, allowedRoles ...types.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
			return
		}

		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrContentNotFound))
			return
		}

		user, err := db.GetUserByUsername(c, username)
		if err != nil {
			ac.Log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.WithError(response.ErrContentNotFound))
			return
		}

		roles, err := db.GetUserRoles(c, user.ID)
		if err != nil {
			ac.Log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.WithError(response.ErrContentNotFound))
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if slices.Contains(roles, role) {
				isAllowed = true
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrInsufficientPermission))
			return
		}

		c.Next()
	}
}

// Authz middleware puts roles into the handler context to make it available for the handlers.
func Authz(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
			return
		}

		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrContentNotFound))
			return
		}

		user, err := db.GetUserByUsername(c, username)
		if err != nil {
			ac.Log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.WithError(response.ErrContentNotFound))
			return
		}

		roles, err := db.GetUserRoles(c, user.ID)
		if err != nil {
			ac.Log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.WithError(response.ErrContentNotFound))
			return
		}

		c.Set(common.RolesCtxKey, roles)
		c.Next()
	}
}
