package router

import (
	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// commitsRoutes : Gestion des commits
func commitsRoutes(e *echo.Echo, g *echo.Group) {
	commitsGroup := g.Group("/commits")

	// Groupe protégé
	commitsGroup.Use(middleware.JWT([]byte(config.JWTSecretKey)))

	// Liste des routes
	commitsGroup.GET("", controller.GetCommitsHandler)
	commitsGroup.GET("/:id", controller.GetCommitHandler)
}
