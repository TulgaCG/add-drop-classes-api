package server

import (
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"gitlab.com/go-box/pongo2gin/v6"

	"github.com/TulgaCG/add-drop-classes-api/pkg/auth"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/lecture"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/role"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
	"github.com/TulgaCG/add-drop-classes-api/pkg/user"
)

func New(db *gendb.Queries, log *slog.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.HTMLRender = pongo2gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", pongo2.Context{})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", pongo2.Context{})
	})

	r.GET("/lectures", middleware.Log(log), middleware.Database(db), middleware.Authentication(db), func(c *gin.Context) {

		username, _ := c.Cookie("username")
		user, _ := db.GetUserCredentialsWithUsername(c, username)

		userLectures, _ := db.GetUserLectures(c, user.ID)

		lectures, err := db.ListLectures(c)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}

		c.HTML(http.StatusOK, "lecture.html", pongo2.Context{
			"lectures":  lectures,
			"lectureExist": func (lec string) bool {
				for _, lecture := range userLectures {
					if lecture.Code == lec {
						return true
					}
				}
				return false
			},
		})
	})

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

	g3 := r.Group("/api", middleware.Log(log), middleware.Database(db), middleware.Authentication(db), middleware.Authorization(db, types.RoleAdmin))
	g3.DELETE("/roles/:uid/:rid", role.RemoveFromUserHandler)
	g3.POST("/roles", role.AddToUserHandler)

	g4 := r.Group("/api", middleware.Log(log), middleware.Database(db), middleware.Authentication(db), middleware.Authorization(db, types.RoleAdmin, types.RoleTeacher))
	g4.GET("/users", user.ListHandler)

	return r
}
