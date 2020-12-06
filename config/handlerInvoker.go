package config

import (
	"github.com/eminetto/clean-architecture-go-v2/api/handler"
	"go.uber.org/fx"
)

// HandlerInvoker used for invoking handler registing
var HandlerInvoker = fx.Options(
	fx.Invoke(handler.MakeBookHandlers),
	fx.Invoke(handler.MakeUserHandlers),
	fx.Invoke(handler.MakeLoanHandlers),
)
