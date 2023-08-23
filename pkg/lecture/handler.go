package lecture

import (
	"fmt"
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

func GetFromUserHandler(c *gin.Context) {
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		log.Error("failed to get roles from gin context")
	}

	roles, ok := r.([]types.Role)
	if !ok {
		log.Error("failed type assertion roles")
	}
	
	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := db.GetUserByUsername(c, username)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrContentNotFound))
			return
		}

		id = int(self.ID)
	}

	lectures, err := getLecturesFromUser(c, db, types.UserID(id))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(lectures))
}

func AddToUserHandler(c *gin.Context) {
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

	var req AddLectureToUserRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		log.Error("failed to get roles from gin context")
	}

	roles, ok := r.([]types.Role)
	if !ok {
		log.Error("failed type assertion roles")
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		req.Username = username
	}

	row, err := addLectureToUser(c, db, req.Username, req.LectureCode)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(fmt.Errorf("failed to add lecture")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func RemoveFromUserHandler(c *gin.Context) {
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

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		log.Error("failed to get roles from gin context")
	}

	roles, ok := r.([]types.Role)
	if !ok {
		log.Error("failed type assertion roles")
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := db.GetUserByUsername(c, username)
		if err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrContentNotFound))
			return
		}

		uid = int(self.ID)
	}

	lid, err := strconv.Atoi(c.Param("lid"))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	err = removeLectureFromUser(c, db, types.UserID(uid), types.LectureID(lid))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusNotAcceptable, response.WithError(fmt.Errorf("lecture not found in user")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(""))
}
