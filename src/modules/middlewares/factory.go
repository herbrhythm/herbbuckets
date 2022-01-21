package middlewares

import (
	"github.com/herb-go/herb-drivers/middleware/hiredmiddleware"
	"github.com/herb-go/herb-drivers/middlewarecondition/requestpatterncondition"
	"github.com/herb-go/herb/middleware/middlewarefactory"
)

func init() {
	//Register time condition factory.
	middlewarefactory.DefaultContext.RegisterConditionFactory("time", middlewarefactory.NewTimeConditionFactory())
	//Register request pattern condition factory
	middlewarefactory.DefaultContext.RegisterConditionFactory("pattern", requestpatterncondition.NewConditionFactory())
	//Register reponse condition factory.
	middlewarefactory.DefaultContext.RegisterFactory("response", middlewarefactory.NewResponseFactory())
	//Register hiredmiddleware factory.
	middlewarefactory.DefaultContext.RegisterFactory("hiredmiddleware", hiredmiddleware.NewFactory())
}
