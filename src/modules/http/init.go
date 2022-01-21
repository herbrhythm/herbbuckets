package http

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/routers"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/util"
	"github.com/herb-go/util/httpserver"
)

//ModuleName module name
const ModuleName = "999http"

//App Main applactions. to serve http
var App = middleware.New()

//Middlewares middlewares which should be used in whole app
var Middlewares = func() middleware.Middlewares {
	return middleware.Middlewares{
		app.HTTP.Forwarded.ServeMiddleware,
		app.HTTP.Hosts.ServeMiddleware,
		app.HTTP.ErrorPages.Middleware(),
		httpserver.RecoverMiddleware(nil),
		app.HTTP.Headers.ServeMiddleware,
	}
}

//Start start app as http server
var Start = func() {
	if app.HTTP.Config.Disabled {
		return
	}
	var Server = app.HTTP.Config.Server()
	httpserver.MustListenAndServeHTTP(Server, &app.HTTP.Config, App)
	util.WaitingQuit()
	defer util.Bye()
	httpserver.ShutdownHTTP(Server)

}

func init() {

	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		App.
			Use(Middlewares()...).
			Handle(routers.RouterFactory.CreateRouter())
		go Start()
	})
}
