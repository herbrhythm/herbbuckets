package overseers

import (
	"net/http"

	overseer "github.com/herb-go/herb-drivers/overseers/middlewareoverseer"
	worker "github.com/herb-go/worker"
)

//MiddlewareWorker empty middleware worker.
var MiddlewareWorker func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

//MiddlewareOverseer middleware overseer
var MiddlewareOverseer = worker.NewOrverseer("middleware", &MiddlewareWorker)

func init() {
	MiddlewareOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(MiddlewareOverseer)
	})
	worker.Appoint(MiddlewareOverseer)
}
