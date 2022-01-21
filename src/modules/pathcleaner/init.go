package pathcleaner

import (
	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "900pathcleaner"

func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.InitOrderByName(ModuleName)
	})
}
