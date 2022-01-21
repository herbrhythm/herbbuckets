package httpinfo

import (
	"herbbuckets/modules/app"

	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "100httpinfo"

func init() {
	util.StageInit.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.Must(app.Httpinfo.Register())
		util.Must(app.Presethttpinfo.Register())
		util.Must(app.Httpinfo.ApplyToFields())
		util.Must(app.Presethttpinfo.ApplyToFields())
		util.InitOrderByName(ModuleName)
	})
}
