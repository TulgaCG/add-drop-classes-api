package lecture

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/domain/lecture"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func GetFromUser(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		ac.Log.Error(response.ErrFailedToFindRolesInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles from context")))
		return
	}

	roles, ok := r.([]types.Role)
	if !ok {
		ac.Log.Error("failed type assertion roles")
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles")))
		return
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := ac.DB.GetUserByUsername(c, username)
		if err != nil {
			ac.Log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrContentNotFound))
			return
		}

		id = int(self.ID)
	}

	lectures, err := lecture.GetLecturesFromUser(c, ac.DB, types.UserID(id))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(lectures))
}

func AddToUser(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	var req AddLectureToUserRequest
	if err := c.BindJSON(&req); err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		ac.Log.Error(response.ErrFailedToFindRolesInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles from context")))
		return
	}

	roles, ok := r.([]types.Role)
	if !ok {
		ac.Log.Error("failed type assertion roles")
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles")))
		return
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		req.Username = username
	}

	row, err := lecture.AddLectureToUser(c, ac.DB, req.Username, req.LectureCode)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to add lecture")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func RemoveFromUser(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		ac.Log.Error(response.ErrFailedToFindRolesInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles from context")))
		return
	}

	roles, ok := r.([]types.Role)
	if !ok {
		ac.Log.Error("failed type assertion roles")
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles")))
		return
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := ac.DB.GetUserByUsername(c, username)
		if err != nil {
			ac.Log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrContentNotFound))
			return
		}

		uid = int(self.ID)
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	err = lecture.RemoveLectureFromUser(c, ac.DB, types.UserID(uid), types.LectureID(id))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusNotAcceptable, response.WithError(errors.New("lecture not found in user")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(""))
}
