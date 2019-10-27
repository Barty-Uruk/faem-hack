package router

import (

	// cfg "gitlab.com/faemproject/backend/hack/services/crm/config"

	echoprometheus "github.com/0neSe7en/echo-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	h "gitlab.com/faemproject/backend/hack/handlers"
	"gitlab.com/faemproject/backend/hack/models"
)

const apiversion = "v2"

// Init - binding middleware and setup routers
func Init(l *logrus.Logger) *echo.Echo {

	e := echo.New()

	// logrus.SetLevel(logrus.DebugLevel)
	// e.Use(logrusmiddleware.Hook())
	newLogger(l)
	e.Use(Hook())

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Metrics
	e.Use(echoprometheus.NewMetric())
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	jwtSec := models.JwtSec
	jwtGroup := e.Group("/api/" + apiversion + "/login")
	jwtGroup.POST("/", h.LoginHandler)
	jwtGroup.POST("/registration", h.CreateUser)
	o := e.Group("/api/" + apiversion)
	o.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(jwtSec),
	}))
	o.GET("/project/:id", h.GetProject)
	o.GET("/project/statuses/needcheck", h.GetStatusesForTeacher)
	o.GET("/project/statuses/:id", h.GetStatusesByProject)
	o.POST("/project/status", h.NewProjectStatus)
	o.PUT("/group/addusers", h.AddUsersToGroup)
	o.GET("/user/:uuid", h.GetUser)
	o.GET("/users", h.GetUsers)
	return e
}
