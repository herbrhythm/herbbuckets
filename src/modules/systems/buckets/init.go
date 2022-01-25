package buckets

import (
	"fmt"
	"herbbuckets/modules/app"

	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "900systems.buckets"

func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.InitOrderByName(ModuleName)
		buckets := app.Buckets
		for _, v := range *buckets {
			name := v.Name
			if Buckets[name] != nil {
				panic(fmt.Errorf("bucket [%s] exists", name))
			}
			b, err := CreateBucket(v)
			if err != nil {
				panic(err)
			}
			Buckets[name] = b
		}
	})
}
