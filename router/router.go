package router

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StartServer : Démarrage du serveur
func StartServer(port int) {
	// Initialisation du serveur
	e := initServer()

	// Version de l'API
	g := e.Group("/v1")

	// Liste des routes
	authRoutes(e, g)
	usersRoutes(e, g)
	commitsRoutes(e, g)

	// Lancement du serveur
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}

// initServer : Initialisation du serveur
func initServer() *echo.Echo {
	e := echo.New()

	// Logger
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} |  ${status} | ${latency_human}\t| ${method}\t${uri}\n",
	}))

	// Recover
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Secure
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// Favicon
	e.File("/favicon.ico", "assets/images/favicon.ico")

	return e
}
