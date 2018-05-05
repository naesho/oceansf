package context

import (
	"github.com/labstack/echo"
	"github.com/ohsean53/oceansf/db"
	"github.com/ohsean53/oceansf/cache"
)

// life cycle = Per Request
type SessionContext struct {
	echo.Context
	DB *db.DB
	Cache *cache.Cache
}

