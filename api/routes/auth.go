package routes

import (
	"github.com/labstack/echo"
	ds "google/datastore"
	"net/http"
	"response"
)

type Auth struct {
	Username string `bind:"required"`
	password string `bind:"required"`
	hashAddress string `bind:"required"`
}

func SignUp (c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// TODO : Validator

	auth := &Auth{
		Username: username,
		password: password,
		hashAddress: "",	// TODO : Create Hash Address And QR code
	}

	key, err := ds.InsertUser(c, auth)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.DatastoreError())
	}
	
	return c.JSON(http.StatusOK, key.IntID())
}

func Login (c echo.Context) error {

	return nil
}