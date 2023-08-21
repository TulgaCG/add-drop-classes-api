package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func CreateHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	var req CreateUserRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	row, err := createUser(c, db, req.Username, req.Password)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(err))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func ListHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	roles := c.GetStringSlice(common.RolesCtxKey)
	if !slices.Contains(roles, "admin") && !slices.Contains(roles, "teacher") {
		log.Error(response.ErrInsufficientPermission.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(response.ErrInsufficientPermission))
		return
	}

	rows, err := listUsers(c, db)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(err))
		return
	}

	c.JSON(http.StatusOK, response.WithData(rows))
}

func GetHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	roles := c.GetStringSlice(common.RolesCtxKey)
	if !slices.Contains(roles, "admin") && !slices.Contains(roles, "teacher") {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := getUserByUsername(c, db, username)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrContentNotFound))
			return
		}

		c.JSON(http.StatusOK, response.WithData(self))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	row, err := getUser(c, db, types.UserID(id))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(err))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func UpdateHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	var req UpdateUserRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	roles := c.GetStringSlice(common.RolesCtxKey)
	if !slices.Contains(roles, "admin") && !slices.Contains(roles, "teacher") {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := updateUser(c, db, gendb.UpdateUserParams{Username: req.NewUsername, Password: req.NewPassword})
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
			return
		}

		c.JSON(http.StatusOK, response.WithData(self))
		return
	}

	u, err := getUserByUsername(c, db, req.Username)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(err))
		return
	}

	row, err := updateUser(c, db, gendb.UpdateUserParams{
		ID:       u.ID,
		Username: req.NewUsername,
		Password: req.NewPassword,
	})
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(err))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}
