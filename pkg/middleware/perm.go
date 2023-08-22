package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
)

func PermMiddleware(db *gendb.Queries, allowedRoles []string, blockHandlers bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
			return
		}

		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrContentNotFound))
			return
		}

		user, err := db.GetUserByUsername(c, username)
		if err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.WithError(response.ErrContentNotFound))
			return
		}

		roles, err := db.GetUserRoles(c, user.ID)
		if err != nil {
			log.Error(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, response.WithError(response.ErrContentNotFound))
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if slices.Contains(roles, role) {
				isAllowed = true
			}
		}

		if blockHandlers && !isAllowed {
			log.Error(response.ErrInsufficientPermission.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrInsufficientPermission))
			return
		}

		c.Set(common.RolesCtxKey, roles)
		c.Next()
	}
}
