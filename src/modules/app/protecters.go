package app

import (
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/herbmodule-drivers/protecter-drivers/protecters/hiredprotecters"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//Protecters config instance of protecter.
var Protecters = &hiredprotecters.ProtectersConfig{}

var syncProtecters atomic.Value

//StoreProtecters atomically store protecter config
func (a *appSync) StoreProtecters(c *hiredprotecters.ProtectersConfig) {
	syncProtecters.Store(c)
}

//LoadProtecters atomically load protecter config
func (a *appSync) LoadProtecters() *hiredprotecters.ProtectersConfig {
	v := syncProtecters.Load()
	if v == nil {
		return nil
	}
	return v.(*hiredprotecters.ProtectersConfig)
}

func init() {
	//Register loader which will be execute when Config.LoadAll func be called.
	//You can put your init code after load.
	//You must panic if any error rasied when init.
	config.RegisterLoader(util.ConfigFile("/protecters.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Protecters))
		Sync.StoreProtecters(Protecters)
	})
}
