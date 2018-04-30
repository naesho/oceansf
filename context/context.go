package context

import (
	"github.com/labstack/echo"
	"github.com/naesho/oceansf/db"
	"github.com/naesho/oceansf/cache"
)

// life cycle = Per Request
type RequestContext struct {
	echo.Context
	DB *db.DB
	Cache *cache.Cache
}

