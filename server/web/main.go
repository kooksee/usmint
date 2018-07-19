package web

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Run(port string) error {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// v1:=e.Group("v1", middleware.BasicAuth(func(s string, s2 string, context echo.Context) (bool, error) {
	// 	return true, nil
	// }))

	// e.Pre(middleware.HTTPSNonWWWRedirect())
	// e.Pre(middleware.WWWRedirect())
	// e.Use(middleware.HTTPSRedirectWithConfig(middleware.RedirectConfig{
	// 	Code: http.StatusTemporaryRedirect,
	// }))
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status}\n",
	// }))

	e.Use(middleware.RequestID())
	e.Validator = &CustomValidator{validator: validator.New()}

	// e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{URL: url1}, {URL: url2}})))

	// 向链上发送tx,需要用户自己传入一个ID
	e.POST("/", txPost)

	// 获得异步发送的结果
	e.GET("/:id", txGet)

	// Start server
	return e.Start(f(":%s", port))
}
