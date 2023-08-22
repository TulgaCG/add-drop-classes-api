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

	g1 := r.Group("/api", middleware.Log(log), middleware.Database(db))
	g1.GET("/logout", auth.LogoutHandler)
	g1.POST("/login", auth.LoginHandler)
	g1.POST("/users", user.CreateHandler)

	g2 := r.Group("/api", middleware.Log(log), middleware.Database(db), middleware.Authentication(db), middleware.HandlerAuthorization(db))
	g2.GET("/users/:id", user.GetHandler)
	g2.GET("/lectures/:id", lecture.GetFromUserHandler)
	g2.POST("/lectures", lecture.AddToUserHandler)
	g2.DELETE("/lectures/:uid/:lid", lecture.RemoveFromUserHandler)
	g2.PUT("/users", user.UpdateHandler)

	g3 := r.Group("/api", middleware.Log(log), middleware.Database(db), middleware.Authentication(db), middleware.Authorization(db, "admin"))
	g3.DELETE("/roles/:uid/:rid", role.RemoveFromUserHandler)
	g3.POST("/roles", role.AddToUserHandler)

	r.GET("/users", user.ListHandler, middleware.Log(log), middleware.Database(db), middleware.Authentication(db), middleware.Authorization(db, "admin", "teacher"))

	return r
}
