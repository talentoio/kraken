package resthttp

import (
	"kraken/internal/services"

	"github.com/labstack/echo/v4"
)

type RouterDependencies struct {
	LTPService services.LTPService
}

//	@title			LTP API
//	@version		1.0
//	@description	Kraken parser LTP API

//	@contact.name	Yuri Gasparyan

//	@host	localhost:8080

// RegisterRoutes create echo router from dependencies.
//
//nolint:funlen // no another solution
func RegisterRoutes(
	dependencies *RouterDependencies,
) *echo.Echo {

	engine := echo.New()
	router := engine.Group("/api/v1")

	ltpHandler := NewLTPHandler(dependencies.LTPService)
	ltpRouter := router.Group("/ltp")
	ltpRouter.GET("", ltpHandler.Get)

	return engine
}
