// Copyright 2019 Profects Group B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
