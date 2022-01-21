package app

import (
	"sync/atomic"
	
	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/uniqueid"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//UniqueID unique id module config
var UniqueID *uniqueid.OptionConfig

var syncUniqueID atomic.Value

//StoreUniqueID atomically store unique id module config
func (a *appSync) StoreUniqueID(c *uniqueid.OptionConfig) {
	syncUniqueID.Store(c)
}

//LoadUniqueID atomically load unique id module config
func (a *appSync) LoadUniqueID() *uniqueid.OptionConfig {
	v := syncUniqueID.Load()
	if v == nil {
		return nil
	}
	return v.(*uniqueid.OptionConfig)
}

func init() {
	config.RegisterLoader(util.ConfigFile("/uniqueid.toml"), func(configpath source.Source) {
		UniqueID = uniqueid.NewOptionConfig()
		util.Must(tomlconfig.Load(configpath, UniqueID))
		Sync.StoreUniqueID(UniqueID)
	})
}
