package controller

import (
	"net/http"
	"time"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// PostAuthLoginHandler : Authentification
func PostAuthLoginHandler(c echo.Context) error {
	// Récupération des variables transmises
	// -------------------------------------
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Vérification en base
	// --------------------
	user, err := model.CheckLogin(username, password)

	if err == nil && user.ID != 0 {
		// Création du token d'authentification
		// ------------------------------------
		token := jwt.New(jwt.SigningMethodHS256)

		// Enregistrement de la revendication
		// ----------------------------------
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["username"] = user.Username
		claims["lastname"] = user.Lastname
		claims["firstname"] = user.Firstname
		claims["exp"] = time.Now().Add(config.JWTExp).Unix()

		// Génération du token encodé et envoi dans la réponse
		// ---------------------------------------------------
		t, err := token.SignedString([]byte(config.JWTSecretKey))

		if err != nil {
			return err
		}

		response := model.GetHTTPResponse(
			http.StatusOK,
			"Success",
			map[string]interface{}{
				"token":     t,
				"id":        user.ID,
				"lastname":  user.Lastname,
				"firstname": user.Firstname,
				"fullname":  user.GetFullname(),
				"createdAt": string(user.CreatedAt),
				"deletedAt": string(user.DeletedAt),
			},
		)

		return c.JSON(response.Code, response)
	}

	return echo.ErrUnauthorized
}

// GetAuthLogoutHandler : Déconnexion
func GetAuthLogoutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Logout")
}
