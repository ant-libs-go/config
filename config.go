/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-04 20:27:56
# File Name: config.go
# Description:
####################################################################### */

package config

import (
	"reflect"
	"sync"
	"time"

	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/config/parser"
)

type item struct {
	cfg  interface{}
	opts *options.Options
}

type Config struct {
	lock   sync.RWMutex
	m      map[string]*item
	opts   *options.Options
	parser parser.Parser
}

func New(parser parser.Parser, opts ...options.Option) (r *Config) {
	options := &options.Options{
		CheckInterval: 10,
		OnChangeFn:    func(cfg interface{}) {},
		OnErrorFn:     func(error) {},
	}
	for _, opt := range opts {
		opt(options)
	}

	r = &Config{
		m:      map[string]*item{},
		opts:   options,
		parser: parser,
	}
	go r.changeChecker()

	if Default == nil {
		Default = r
	}
	return
}

func (this *Config) Load(cfg interface{}, opts ...options.Option) (r *Config, err error) {
	name := reflect.TypeOf(cfg).String()
	options := &options.Options{
		Sources:    this.opts.Sources,
		OnChangeFn: this.opts.OnChangeFn,
	}
	for _, opt := range opts {
		opt(options)
	}

	this.lock.Lock()
	this.m[name] = &item{cfg: cfg, opts: options}
	err = this.load(this.m[name])
	this.lock.Unlock()

	r = this
	return
}

// must be locked at call the func
func (this *Config) load(item *item) (err error) {
	ty := reflect.TypeOf(item.cfg)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}
	item.cfg = reflect.New(ty).Interface()

	err = this.parser.Unmarshal(item.cfg, item.opts)
	if err != nil {
		return
	}

	item.opts.OnChangeFn(item.cfg)
	return
}

func (this *Config) changeChecker() {
	var err error
	ticker := time.NewTicker(time.Second * time.Duration(this.opts.CheckInterval))
	for _ = range ticker.C {
		var lastModTime int64
		lastModTime, err = this.parser.GetLastModTime(this.opts)
		if err != nil {
			this.doError(err)
			continue
		}

		if time.Now().Unix()-lastModTime > this.opts.CheckInterval {
			continue
		}

		this.lock.Lock()
		for _, item := range this.m {
			this.doError(this.load(item))
		}
		this.lock.Unlock()
	}
	return
}

func (this *Config) Get(cfg interface{}) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if v, ok := this.m[reflect.TypeOf(cfg).String()]; ok {
		return v.cfg
	}
	return nil
}

func (this *Config) doError(err error) {
	if err == nil {
		return
	}
	this.opts.OnErrorFn(err)
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
