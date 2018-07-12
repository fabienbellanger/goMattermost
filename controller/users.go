package controller

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fabienbellanger/goMattermost/model"
	"github.com/labstack/echo"
)

// userType
type userType struct {
	Username  string  `json:"username" xml:"username"`
	Lastname  string  `json:"lastname" xml:"lastname"`
	Firstname string  `json:"firstname" xml:"firstname"`
	Exp       float64 `json:"exp" xml:"exp"`
}

// GetUsersInfosHandler : Information utilisateur connecté
func GetUsersInfosHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	lastname := claims["lastname"].(string)
	firstname := claims["firstname"].(string)
	exp := claims["exp"].(float64)

	userType := &userType{
		Username:  username,
		Lastname:  lastname,
		Firstname: firstname,
		Exp:       exp,
	}

	response := model.GetHTTPResponse(http.StatusOK, "Success", userType)

	return c.JSON(response.Code, response)
}
