/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:37:24
# File Name: options.go
# Description:
####################################################################### */

package options

type Options struct {
	Sources       []string          // config source
	CheckInterval int64             // file update check interval
	OnChangeFn    func(interface{}) // call it when the file is modified
	OnErrorFn     func(error)       // call it when an error occurs
}

type Option func(o *Options)

func WithCfgSource(inp ...string) Option {
	return func(o *Options) {
		if o.Sources == nil {
			o.Sources = []string{}
		}
		o.Sources = append(o.Sources, inp...)
	}
}

func WithCheckInterval(inp int64) Option {
	return func(o *Options) {
		o.CheckInterval = inp
	}
}

func WithOnErrorFn(inp func(error)) Option {
	return func(o *Options) {
		o.OnErrorFn = inp
	}
}

func WithOnChangeFn(inp func(cfg interface{})) Option {
	return func(o *Options) {
		o.OnChangeFn = inp
	}
}

type OpOptions struct {
	OnChangeFn func(interface{}) // call it when the file is modified
	OnErrorFn  func(error)       // call it when an error occurs
}

type OpOption func(o *OpOptions)

func WithOpOnErrorFn(inp func(error)) OpOption {
	return func(o *OpOptions) {
		o.OnErrorFn = inp
	}
}

func WithOpOnChangeFn(inp func(cfg interface{})) OpOption {
	return func(o *OpOptions) {
		o.OnChangeFn = inp
	}
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
