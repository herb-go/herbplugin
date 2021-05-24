package herbplugin

import "context"

type ContextName string

const ContextNameFinished = ContextName("finished")

func GetFinished(ctx context.Context) bool {
	result := ctx.Value(ContextNameFinished).(*bool)
	return *result
}

var Finish = func(ctx context.Context, plugin Plugin) {
	result := ctx.Value(ContextNameFinished).(*bool)
	*result = true
}

type Process func(ctx context.Context, plugin Plugin, next func(ctx context.Context, plugin Plugin))

func ComposeProcess(series ...Process) Process {
	return func(ctx context.Context, plugin Plugin, receiver func(ctx context.Context, plugin Plugin)) {
		if len(series) == 0 {
			receiver(ctx, plugin)
			return
		}
		series[0](ctx, plugin, func(newctx context.Context, plugin Plugin) {
			ComposeProcess(series[1:]...)(newctx, plugin, receiver)
		})
	}
}

func Exec(plugin Plugin, p ...Process) bool {
	finished := false
	ctx := context.WithValue(context.Background(), ContextNameFinished, &finished)
	ComposeProcess(p...)(ctx, plugin, Finish)
	return finished
}
