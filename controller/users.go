package controller

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// userType
type userType struct {
	Username string  `json:"username" xml:"username"`
	Exp      float64 `json:"exp" xml:"exp"`
}

// GetUsersInfosHandler : Information utilisateur connecté
func GetUsersInfosHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	exp := claims["exp"].(float64)

	fmt.Println(claims)
	u := &userType{
		Username: username,
		Exp:      exp,
	}

	return c.JSON(http.StatusOK, u)
}
