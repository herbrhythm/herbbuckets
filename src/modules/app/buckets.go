package app

import (
	"herbbuckets/modules/app/bucketconfig"
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//BucketsConfig buckets config data struct.
//Struct must  unmarshaleable by Toml lib.
//You should comment this struct if you use third party config struct.
type BucketsConfig struct {
	Buckets []*bucketconfig.Config
}

//Buckets config instance of buckets.
var Buckets = &BucketsConfig{}

var syncBuckets atomic.Value

//StoreBuckets atomically store buckets config
func (a *appSync) StoreBuckets(c *BucketsConfig) {
	syncBuckets.Store(c)
}

//LoadBuckets atomically load buckets config
func (a *appSync) LoadBuckets() *BucketsConfig {
	v := syncBuckets.Load()
	if v == nil {
		return nil
	}
	return v.(*BucketsConfig)
}

func init() {
	//Register loader which will be execute when Config.LoadAll func be called.
	//You can put your init code after load.
	//You must panic if any error rasied when init.
	config.RegisterLoader(util.ConfigFile("/buckets.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Buckets))
		Sync.StoreBuckets(Buckets)
	})
}
