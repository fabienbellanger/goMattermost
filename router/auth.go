package router

import (
	"github.com/fabienbellanger/goMattermost/controller"
	"github.com/labstack/echo"
)

// authRoutes : Partie authentification
func authRoutes(e *echo.Echo, g *echo.Group) {
	authGroup := g.Group("/auth")

	// Liste des routes
	authGroup.POST("/login", controller.PostAuthLoginHandler)
}
