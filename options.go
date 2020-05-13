package apollo

import (
	"context"
	"github.com/micro/go-micro/v2/config/source"
)

type apolloConfPath struct{}
type namespaceName struct{}

func WithApolloConfPath(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, apolloConfPath{}, c)
	}
}

func WithNamespaceName(c string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, namespaceName{}, c)
	}
}