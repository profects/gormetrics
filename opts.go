package gormetrics

type RegisterOpt func(o *pluginOpts)

type pluginOpts struct {
	prometheusNamespace string
	gormPluginScope     string
}

// WithPrometheusNamespace sets a different namespace for the exported metrics.
// The default namespace is "gormetrics".
func WithPrometheusNamespace(ns string) RegisterOpt {
	return func(o *pluginOpts) {
		o.prometheusNamespace = ns
	}
}

// WithGORMPluginScope sets a different plugin scope for the configured callbacks.
// The default plugin scope is "gormetrics".
func WithGORMPluginScope(pluginScope string) RegisterOpt {
	return func(o *pluginOpts) {
		o.gormPluginScope = pluginScope
	}
}

// defaultCallbackHandlerOpts creates a new pluginOpts instance with the default values.
func defaultCallbackHandlerOpts() *pluginOpts {
	return &pluginOpts{
		prometheusNamespace: "gormetrics",
		gormPluginScope:     "gormetrics",
	}
}

// getOpts creates a pluginOpts instance based on multiple user-defined options based on the default
// options. See defaultCallbackOpts for the default options.
func getOpts(opts []RegisterOpt) *pluginOpts {
	c := defaultCallbackHandlerOpts()
	for _, o := range opts {
		o(c)
	}
	return c
}
