/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:37:24
# File Name: options.go
# Description:
####################################################################### */

package options

type Options struct {
	Files         []string          // config main file
	CheckInterval int64             // file update check interval
	OnChangeFn    func(interface{}) // call it when the file is modified
	OnErrorFn     func(error)       // call it when an error occurs
}

type Option func(o *Options)

func WithCfgFile(inp string) Option {
	return func(o *Options) {
		if o.Files == nil {
			o.Files = []string{}
		}
		o.Files = append(o.Files, inp)
	}
}

func WithCfgFiles(inp ...string) Option {
	return func(o *Options) {
		if o.Files == nil {
			o.Files = []string{}
		}
		o.Files = append(o.Files, inp...)
	}
}

func WithCheckInterval(inp int64) Option {
	return func(o *Options) {
		o.CheckInterval = inp
	}
}

func WithOnChangeFn(inp func(interface{})) Option {
	return func(o *Options) {
		o.OnChangeFn = inp
	}
}

func WithOnErrorFn(inp func(error)) Option {
	return func(o *Options) {
		o.OnErrorFn = inp
	}
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
