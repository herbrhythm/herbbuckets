package routers

import (
	"herbbuckets/modules/middlewares"
	"herbbuckets/modules/protecters"
	bucketsactions "herbbuckets/modules/systems/buckets/actions"
	bucketsmiddlewares "herbbuckets/modules/systems/buckets/middlewares"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/errorpage"
	"github.com/herb-go/herb/middleware/misc"
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
	Router.GET("/current").Handle(bucketsactions.ActionCurrent)
	Router.StripPrefix("/grantdownloadurl").
		Use(
			misc.MethodGET.ServeMiddleware,
			bucketsmiddlewares.MiddlewarePath,
			bucketsmiddlewares.MiddlewareAuthViewBucket,
		).Handle(bucketsactions.ActionGrantDownloadURL)
	Router.StripPrefix("/grantuploadurl").
		Use(
			misc.MethodGET.ServeMiddleware,
			bucketsmiddlewares.MiddlewareBase,
			bucketsmiddlewares.MiddlewareAuthUploadBucket,
		).Handle(bucketsactions.ActionGrantUploadURL)
	Router.StripPrefix("/info").
		Use(
			misc.MethodGET.ServeMiddleware,
			bucketsmiddlewares.MiddlewarePath,
			bucketsmiddlewares.MiddlewareAuthViewBucket,
		).Handle(bucketsactions.ActionInfo)
	Router.StripPrefix("/content").
		Use(
			misc.MethodGET.ServeMiddleware,
			bucketsmiddlewares.MiddlewarePath,
			bucketsmiddlewares.MiddlewareAuthViewBucket,
		).Handle(bucketsactions.ActionContent)
	Router.StripPrefix("/remove").
		Use(
			misc.MethodPOST.ServeMiddleware,
			bucketsmiddlewares.MiddlewarePath,
			bucketsmiddlewares.MiddlewareAuthUploadBucket,
		).Handle(bucketsactions.ActionRemove)

	Router.StripPrefix("/save").
		Use(
			misc.MethodPOST.ServeMiddleware,
			bucketsmiddlewares.MiddlewareBase,
			bucketsmiddlewares.MiddlewareAuthUploadBucket,
		).Handle(bucketsactions.ActionSave)

	return Router
})
