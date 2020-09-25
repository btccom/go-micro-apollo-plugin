package apollo

import (
	"context"
	"github.com/micro/go-micro/v2/config/source"
)

type apolloConfPath struct{}
type namespaceName struct{}
type isCustomConfig struct{}
type customConfig struct{}

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

func WithIsCustomConfig(c bool) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, isCustomConfig{}, c)
	}
}

func WithCustomConfig(c *CustomConfig) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, customConfig{}, c)
	}
}
