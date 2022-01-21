package overseers

import (
	"github.com/herb-go/herb/middleware/action"
	overseer "github.com/herb-go/herb-drivers/overseers/actionoverseer"
	worker "github.com/herb-go/worker"
)

//ActionWorker empty cache worker.
var ActionWorker = action.New(nil)

//ActionOverseer cache overseer
var ActionOverseer = worker.NewOrverseer("action", &ActionWorker)

func init() {
	ActionOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(ActionOverseer)
	})
	worker.Appoint(ActionOverseer)
}
