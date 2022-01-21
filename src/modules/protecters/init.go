package protecters

import (
	"herbbuckets/modules/app"
	"net/http"

	"github.com/herb-go/herbmodules/protecter/protectermanager"

	"github.com/herb-go/herbmodules/protecter"

	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "200protecters"

//ProtectMiddleware protect middleware
func ProtectMiddleware(name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return protecter.ProtectMiddleware(protectermanager.Register(name))
}

//ProtectWith protect given handler with Protecter
func ProtectWith(name string, h http.Handler) http.Handler {
	return protecter.ProtectWith(protectermanager.Register(name), h)
}

//Get get protecter by given name.
func Get(name string) *protecter.Protecter {
	return protectermanager.Register(name)
}

func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		protectermanager.Debug = app.Development.Debug
		util.Must(app.Protecters.Apply())
		util.InitOrderByName(ModuleName)
	})
}
