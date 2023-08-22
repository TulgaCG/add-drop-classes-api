package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/auth"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/lecture"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/role"
	"github.com/TulgaCG/add-drop-classes-api/pkg/user"
)

func New(db *gendb.Queries, log *slog.Logger) *gin.Engine {
	r := gin.Default()

	g1 := r.Group("/api", middleware.LogMiddleware(log), middleware.DBMiddleware(db))
	g1.GET("/logout", auth.LogoutHandler)
	g1.POST("/login", auth.LoginHandler)
	g1.POST("/users", user.CreateHandler)

	g2 := r.Group("/api", middleware.LogMiddleware(log), middleware.DBMiddleware(db), middleware.Authentication(db))
	g2.GET("/users", user.ListHandler, middleware.Authorization(db, "admin", "teacher"))
	g2.GET("/users/:id", user.GetHandler, middleware.HandlerAuthorization(db))
	g2.GET("/lectures/:id", lecture.GetFromUserHandler, middleware.HandlerAuthorization(db))
	g2.POST("/roles", role.AddToUserHandler, middleware.Authorization(db, "admin"))
	g2.POST("/lectures", lecture.AddToUserHandler, middleware.HandlerAuthorization(db))
	g2.DELETE("/lectures/:uid/:lid", lecture.RemoveFromUserHandler, middleware.HandlerAuthorization(db))
	g2.DELETE("/roles/:uid/:rid", role.RemoveFromUserHandler, middleware.Authorization(db, "admin"))
	g2.PUT("/users", user.UpdateHandler, middleware.HandlerAuthorization(db))

	return r
}
