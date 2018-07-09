package router

import (
	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// usersRoutes : Gestion des utilisateurs
func usersRoutes(e *echo.Echo, g *echo.Group) {
	usersGroup := g.Group("/users")

	// Groupe protégé
	usersGroup.Use(middleware.JWT([]byte(config.JWTSecretKey)))

	// Liste des routes
	usersGroup.GET("/infos", controller.GetUsersInfosHandler)
}
