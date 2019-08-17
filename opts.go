package gormetrics

// RegisterOpt if a function that operates on pluginOpts, configuring one or
// more parameters of the plugin options.
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
func WithGORMPluginScope(s string) RegisterOpt {
	return func(o *pluginOpts) {
		o.gormPluginScope = s
	}
}

// defaultPluginOpts creates a new pluginOpts instance with the default values.
func defaultPluginOpts() *pluginOpts {
	return &pluginOpts{
		prometheusNamespace: "gormetrics",
		gormPluginScope:     "gormetrics",
	}
}

// getOpts creates a pluginOpts instance based on multiple user-defined options based
// on the default options. See defaultPluginOpts for the default options.
func getOpts(opts []RegisterOpt) *pluginOpts {
	c := defaultPluginOpts()
	for _, o := range opts {
		o(c)
	}
	return c
}
