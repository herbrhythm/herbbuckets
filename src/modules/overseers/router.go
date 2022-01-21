package overseers

import (
	overseer "github.com/herb-go/herb-drivers/overseers/routeroverseer"
	"github.com/herb-go/herb/middleware/router"
	worker "github.com/herb-go/worker"
)

//RouterFactoryWorker empty router factory worker.
var RouterFactoryWorker *router.Factory

//RouterFactoryOverseer router overseer
var RouterFactoryOverseer = worker.NewOrverseer("router", &RouterFactoryWorker)

func init() {
	RouterFactoryOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(RouterFactoryOverseer)
	})
	worker.Appoint(RouterFactoryOverseer)
}
