package routers

import (
	"herbbuckets/modules/app"
	bucketsactions "herbbuckets/modules/systems/buckets/actions"
	bucketsmiddlewares "herbbuckets/modules/systems/buckets/middlewares"

	//"herbbuckets/modules/actions"
	"github.com/herb-go/herb/file/simplehttpserver"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/middleware/router/httprouter"
	"github.com/herb-go/util"
)

//RouterFactory base router factory.
var RouterFactory = router.NewFactory(func() router.Router {
	var Router = httprouter.New()

	//Only host assests folder if folder exisits
	if app.Assets.URLPrefix != "" {
		Router.StripPrefix(app.Assets.URLPrefix).
			Use(AssestsMiddlewares()...).
			HandleFunc(simplehttpserver.ServeFolder(util.Resources(app.Assets.Location)))
	}
	Router.StripPrefix("/api").
		Use(APIMiddlewares()...).
		Handle(RouterAPIFactory.CreateRouter())
	Router.StripPrefix("/file").
		Use(bucketsmiddlewares.MiddlewarePath).
		Handle(bucketsactions.ActionDownload)
	//var RouterHTML = newHTMLRouter()
	//Router.StripPrefix("/page").Use(HTMLMiddlewares()...).Handle(RouterHTML)

	//Router.GET("/").Use(HTMLMiddlewares()...).Handle(actions.IndexAction)

	return Router
})
