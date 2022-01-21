package api

import (
	"github.com/herb-go/util"
)
//ModuleName module name
const ModuleName="900systems.api"

func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.InitOrderByName(ModuleName)
	})
}
