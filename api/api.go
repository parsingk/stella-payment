package api

import (
	"github.com/labstack/echo"
	"middlewares"
	"net/http"
	"gopkg.in/go-playground/validator.v9"
	"api/routes"
)


type CustomValidator struct {
	validator *validator.Validate
}

var Server = createMux()

func createMux() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world")
	})

	auth := e.Group("/auth")

		auth.POST("/sign-up", routes.SignUp)
		auth.PUT("/login", routes.Login)

	user := e.Group("/user")

		user.GET("/:userId", routes.GetInfo)
		user.PUT("/:userId", routes.UpdateInfo)

	payment := e.Group("/payment")

		payment.POST("/send", routes.Send)

	return e
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

