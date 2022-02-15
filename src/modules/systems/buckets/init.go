package buckets

import (
	"fmt"
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"

	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "900systems.buckets"

func init() {
	util.RegisterModule(ModuleName, func() {
		//Init registered initator which registered by RegisterInitiator
		//util.RegisterInitiator(ModuleName, "func", func(){})
		util.InitOrderByName(ModuleName)
		util.RegisterDataFolder(bucket.BucketsFolder)
		buckets := app.Buckets.Buckets
		for _, v := range buckets {
			name := v.Name
			if !bucket.NameRegexp.MatchString(v.Name) {
				panic(fmt.Errorf("bucket name format error [%s]", name))
			}
			if Buckets[name] != nil {
				panic(fmt.Errorf("bucket [%s] exists", name))
			}
			b, err := CreateBucket(v)
			if err != nil {
				panic(err)
			}
			if !b.Disabled {
				Buckets[name] = b
				util.Must(b.Engine.Start())
				go util.OnQuitAndLogError(b.Engine.Stop)
			}
		}
	})
}
