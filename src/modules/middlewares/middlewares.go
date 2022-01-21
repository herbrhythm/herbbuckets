package middlewares

import (
	"github.com/herb-go/util"
)

//ModuleName module name used in initing and debuging.
const ModuleName = "200Middleware"


func init() {
	util.RegisterModule(ModuleName, func() {
		util.InitOrderByName(ModuleName)
	})
}
