package routers

import (
	"herbbuckets/modules/middlewares"
	bucketsmiddleware "herbbuckets/modules/systems/buckets/middleware"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/errorpage"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/middleware/router/httprouter"
)

//APIMiddlewares middlewares that should used in api requests
var APIMiddlewares = func() middleware.Middlewares {
	return middleware.Middlewares{
		middlewares.MiddlewareCsrfVerifyHeader,
		errorpage.MiddlewareDisable,
		bucketsmiddleware.MiddlewareQuery,
	}
}

//RouterAPIFactory api router factory.
var RouterAPIFactory = router.NewFactory(func() router.Router {
	var Router = httprouter.New()
	//Put your router configure code here
	return Router
})
