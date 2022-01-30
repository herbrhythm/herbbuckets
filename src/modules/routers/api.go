package routers

import (
	"herbbuckets/modules/middlewares"
	"herbbuckets/modules/protecters"
	bucketsactions "herbbuckets/modules/systems/buckets/actions"
	bucketsmiddlewares "herbbuckets/modules/systems/buckets/middlewares"

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
		protecters.ProtectMiddleware("appkey"),
	}
}

//RouterAPIFactory api router factory.
var RouterAPIFactory = router.NewFactory(func() router.Router {
	var Router = httprouter.New()
	//Put your router configure code here
	Router.StripPrefix("/grantdownloadurl").
		Use(
			bucketsmiddlewares.MiddlewarePath,
			bucketsmiddlewares.MiddlewareAuthViewBucket,
		).Handle(bucketsactions.ActionGrantDownloadURL)
	Router.StripPrefix("/grantuploadurl").
		Use(
			bucketsmiddlewares.MiddlewareBase,
			bucketsmiddlewares.MiddlewareAuthUploadBucket,
		).Handle(bucketsactions.ActionGrantUploadURL)
	return Router
})
