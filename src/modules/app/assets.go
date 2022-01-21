package app

import (
	"sync/atomic"

	"github.com/herb-go/herb/file/store"
	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//Assets assets file store.
var Assets = &store.Assets{}

var syncAssets atomic.Value

//StoreAssets atomically store assets config
func (a *appSync) StoreAssets(c *store.Assets) {
	syncAssets.Store(c)
}

//LoadAssets atomically load assets config
func (a *appSync) LoadAssets() *store.Assets {
	v := syncAssets.Load()
	if v == nil {
		return nil
	}
	return v.(*store.Assets)
}

func init() {
	config.RegisterLoader(util.ConstantsFile("/assets.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Assets))
		Sync.StoreAssets(Assets)
	})
}
