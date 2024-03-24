package fxlog

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

type ZerologAdapter struct {
	zerolog.Logger
}

func NewZerologAdapter(logger zerolog.Logger) ZerologAdapter {
	return ZerologAdapter{Logger: logger}
}

func (l ZerologAdapter) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.onStartExecuting(e)
	case *fxevent.OnStartExecuted:
		l.onStartExecuted(e)
	case *fxevent.OnStopExecuting:
		l.onStopExecuting(e)
	case *fxevent.OnStopExecuted:
		l.onStopExecuted(e)
	case *fxevent.Supplied:
		l.supplied(e)
	case *fxevent.Provided:
		l.provided(e)
	case *fxevent.Replaced:
		l.replaced(e)
	case *fxevent.Decorated:
		l.decorated(e)
	case *fxevent.Invoking:
		l.invoking(e)
	case *fxevent.Invoked:
		l.invoked(e)
	case *fxevent.Stopping:
		l.stopping(e)
	case *fxevent.Stopped:
		l.stopped(e)
	case *fxevent.RollingBack:
		l.rollingBack(e)
	case *fxevent.RolledBack:
		l.rolledBack(e)
	case *fxevent.Started:
		l.started(e)
	case *fxevent.LoggerInitialized:
		l.loggerInitialized(e)
	}
}

func (l ZerologAdapter) onStartExecuting(e *fxevent.OnStartExecuting) {
	l.Info().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("OnStart hook executing")
}

func (l ZerologAdapter) onStartExecuted(e *fxevent.OnStartExecuted) {
	if e.Err != nil {
		l.Err(e.Err).Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("OnStart hook failed")
	} else {
		l.Info().Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Str("runtime", e.Runtime.String()).
			Msg("OnStart hook executed")
	}
}

func (l ZerologAdapter) onStopExecuting(e *fxevent.OnStopExecuting) {
	l.Info().Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("OnStop hook executing")
}

func (l ZerologAdapter) onStopExecuted(e *fxevent.OnStopExecuted) {
	if e.Err != nil {
		l.Err(e.Err).Str("callee", e.FunctionName).Str("caller", e.CallerName).Msg("OnStop hook failed")
	} else {
		l.Info().Str("callee", e.FunctionName).Str("caller", e.CallerName).Str("runtime", e.Runtime.String()).
			Msg("OnStop hook executed")
	}
}

func (l ZerologAdapter) supplied(e *fxevent.Supplied) {
	if e.Err != nil {
		l.Err(e.Err).Str("type", e.TypeName).Str("module", e.ModuleName).
			Msg("error encountered while applying options")
	} else {
		l.Info().Str("type", e.TypeName).Str("module", e.ModuleName).Msg("supplied")
	}
}

func (l ZerologAdapter) provided(e *fxevent.Provided) {
	for _, typ := range e.OutputTypeNames {
		l.Info().Str("constructor", e.ConstructorName).
			Str("module", e.ModuleName).
			Str("type", typ).
			Bool("private", e.Private).
			Msg("provided")
	}
	if e.Err != nil {
		l.Err(e.Err).Str("module", e.ModuleName).Msg("error encountered while applying options")
	}
}

func (l ZerologAdapter) replaced(e *fxevent.Replaced) {
	for _, typ := range e.OutputTypeNames {
		l.Info().Str("module", e.ModuleName).Str("type", typ).Msg("replaced")
	}
	if e.Err != nil {
		l.Err(e.Err).Str("module", e.ModuleName).Msg("error encountered while replacing")
	}
}

func (l ZerologAdapter) decorated(e *fxevent.Decorated) {
	for _, typ := range e.OutputTypeNames {
		l.Info().Str("decorator", e.DecoratorName).Str("module", e.ModuleName).Str("type", typ).Msg("decorated")
	}
	if e.Err != nil {
		l.Err(e.Err).Str("module", e.ModuleName).Msg("error encountered while applying options")
	}
}

func (l ZerologAdapter) invoking(e *fxevent.Invoking) {
	l.Info().Str("function", e.FunctionName).Str("module", e.ModuleName).Msg("invoking")
}

func (l ZerologAdapter) invoked(e *fxevent.Invoked) {
	if e.Err != nil {
		l.Err(e.Err).Str("function", e.FunctionName).Str("module", e.ModuleName).Str("stack", e.Trace).
			Msg("invoke failed")
	}
}

func (l ZerologAdapter) stopping(e *fxevent.Stopping) {
	l.Info().Str("signal", strings.ToUpper(e.Signal.String())).
		Msg("received signal")
}

func (l ZerologAdapter) stopped(e *fxevent.Stopped) {
	if e.Err != nil {
		l.Err(e.Err).Msg("stop failed")
	}
}

func (l ZerologAdapter) rollingBack(e *fxevent.RollingBack) {
	l.Err(e.StartErr).Msg("start failed, rolling back")
}

func (l ZerologAdapter) rolledBack(e *fxevent.RolledBack) {
	if e.Err != nil {
		l.Err(e.Err).Msg("rollback failed")
	}
}

func (l ZerologAdapter) started(e *fxevent.Started) {
	if e.Err != nil {
		l.Err(e.Err).Msg("start failed")
	} else {
		l.Info().Msg("started")
	}
}

func (l ZerologAdapter) loggerInitialized(e *fxevent.LoggerInitialized) {
	if e.Err != nil {
		l.Err(e.Err).Msg("custom logger initialization failed")
	} else {
		l.Info().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
	}
}
