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

	g2 := r.Group("/api", middleware.LogMiddleware(log), middleware.DBMiddleware(db), middleware.AuthMiddleware(db))
	g2.GET("/users", user.ListHandler, middleware.PermMiddleware(db, []string{"admin", "teacher"}, true))
	g2.GET("/users/:id", user.GetHandler, middleware.PermMiddleware(db, []string{"admin", "teacher"}, false))
	g2.GET("/lectures/:id", lecture.GetFromUserHandler, middleware.PermMiddleware(db, []string{"admin", "teacher"}, false))
	g2.POST("/roles", role.AddToUserHandler, middleware.PermMiddleware(db, []string{"admin"}, true))
	g2.POST("/lectures", lecture.AddToUserHandler, middleware.PermMiddleware(db, []string{"admin", "teacher"}, false))
	g2.DELETE("/lectures/:uid/:lid", lecture.RemoveFromUserHandler, middleware.PermMiddleware(db, []string{"admin", "teacher"}, false))
	g2.DELETE("/roles/:uid/:rid", role.RemoveFromUserHandler, middleware.PermMiddleware(db, []string{"admin"}, true))
	g2.PUT("/users", user.UpdateHandler, middleware.PermMiddleware(db, []string{"admin", "teacher"}, false))

	return r
}
