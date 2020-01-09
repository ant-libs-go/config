/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-04 20:27:56
# File Name: config.go
# Description:
####################################################################### */

package config

import (
	"sync"
	"time"

	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/config/parser"
	"github.com/ant-libs-go/util"
)

var lock sync.RWMutex

type Config struct {
	cfg    interface{}
	rawCfg interface{}
	opts   *options.Options
	parser parser.Parser
}

func New(cfg interface{}, parser parser.Parser, opts ...options.Option) (r *Config, err error) {
	options := &options.Options{
		CheckInterval: 10,
		OnChangeFn:    func(interface{}) {},
		OnErrorFn:     func(error) {},
	}
	for _, opt := range opts {
		opt(options)
	}

	r = &Config{
		rawCfg: cfg,
		opts:   options,
		parser: parser,
	}
	r.onError(r.Load())
	go r.ChangeChecker()
	return
}

func (this *Config) Load() (err error) {
	cfg := this.rawCfg
	util.DeepCopy(this.rawCfg, this.rawCfg)

	err = this.parser.Unmarshal(cfg, this.opts)
	if err != nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	this.cfg = cfg
	return
}

func (this *Config) ChangeChecker() {
	var err error
	ticker := time.NewTicker(time.Second * time.Duration(this.opts.CheckInterval))
	for _ = range ticker.C {
		var lastModTime int64
		lastModTime, err = this.parser.GetLastModTime(this.opts)
		if err != nil {
			this.onError(err)
			continue
		}

		if time.Now().Unix()-lastModTime <= this.opts.CheckInterval {
			this.onChange()
		}
	}
	return
}

func (this *Config) Get() interface{} {
	lock.RLock()
	defer lock.RUnlock()
	return this.cfg
}

func (this *Config) onChange() {
	err := this.Load()
	if err != nil {
		this.onError(err)
		return
	}
	this.opts.OnChangeFn(this.Get())
}

func (this *Config) onError(err error) {
	if err == nil {
		return
	}
	this.opts.OnErrorFn(err)
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
