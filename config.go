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
	"github.com/ant-libs-go/util"
)

var Instance *Config

type item struct {
	m    interface{}
	hash string
	elem reflect.Value
}

type Config struct {
	lock   sync.RWMutex
	opts   *options.Options
	items  map[string]*item
	parser parser.Parser
}

func NewConfig(parser parser.Parser, opts ...options.Option) (r *Config) {
	options := &options.Options{
		OnChangeFn: func(cfg interface{}) {},
		OnErrorFn:  func(error) {}}
	for _, opt := range opts {
		opt(options)
	}

	Instance = &Config{
		opts:   options,
		items:  map[string]*item{},
		parser: parser}

	go Instance.changeChecker()

	return Instance
}

func (this *Config) changeChecker() {
	if this.opts.CheckInterval == 0 {
		return
	}

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
		for pointer, one := range this.items {
			if err := this.parser.Unmarshal(one.m, this.opts); err != nil {
				this.doError(err)
				continue
			}

			b, _ := util.GobEncode(one.m)
			this.items[pointer] = &item{
				m:    one.m,
				hash: util.Md5String(string(b)),
				elem: reflect.ValueOf(one.m).Elem()}

			if one.hash != this.items[pointer].hash {
				Instance.opts.OnChangeFn(this.items[pointer].m)
			}
		}
		this.lock.Unlock()
	}
	return
}

func (this *Config) doError(err error) {
	if err == nil {
		return
	}
	this.opts.OnErrorFn(err)
}

func Get(cfg interface{}) interface{} {
	Instance.lock.Lock()
	defer Instance.lock.Unlock()

	pointer := reflect.TypeOf(cfg).String()
	if _, ok := Instance.items[pointer]; !ok {
		if err := Instance.parser.Unmarshal(cfg, Instance.opts); err != nil {
			Instance.doError(err)
			return nil
		}

		b, _ := util.GobEncode(cfg)
		Instance.items[pointer] = &item{
			m:    cfg,
			hash: util.Md5String(string(b)),
			elem: reflect.ValueOf(cfg).Elem()}
	}

	if v := Instance.items[pointer]; v != nil {
		reflect.ValueOf(cfg).Elem().Set(v.elem)
	}
	return cfg
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
