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
	m          interface{}
	hash       string
	elem       reflect.Value
	onChangeFn func(interface{})
	onErrorFn  func(error)
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
			this.doError(err, "")
			continue
		}

		if time.Now().Unix()-lastModTime > this.opts.CheckInterval {
			continue
		}

		this.lock.Lock()
		for pointer, one := range this.items {
			if err := this.parser.Unmarshal(one.m, this.opts); err != nil {
				this.doError(err, pointer)
				continue
			}

			oldHash := one.hash

			b, _ := util.JsonEncode(one.m)
			this.items[pointer].m = one.m
			this.items[pointer].hash = util.Md5String(string(b))
			this.items[pointer].elem = reflect.ValueOf(one.m).Elem()

			if oldHash != this.items[pointer].hash {
				this.opts.OnChangeFn(this.items[pointer].m)
				this.items[pointer].onChangeFn(this.items[pointer].m)
			}
		}
		this.lock.Unlock()
	}
	return
}

func (this *Config) doError(err error, pointer string) {
	if err == nil {
		return
	}
	this.opts.OnErrorFn(err)

	if _, ok := this.items[pointer]; ok {
		this.items[pointer].onErrorFn(err)
	}
}

func Get(cfg interface{}, opts ...options.OpOption) interface{} {
	Instance.lock.Lock()
	defer Instance.lock.Unlock()

	pointer := reflect.TypeOf(cfg).String()
	if _, ok := Instance.items[pointer]; !ok {
		options := &options.OpOptions{
			OnChangeFn: func(cfg interface{}) {},
			OnErrorFn:  func(error) {}}
		for _, opt := range opts {
			opt(options)
		}

		Instance.items[pointer] = &item{
			onChangeFn: options.OnChangeFn,
			onErrorFn:  options.OnErrorFn}

		if err := Instance.parser.Unmarshal(cfg, Instance.opts); err != nil {
			Instance.doError(err, pointer)
			return nil
		}

		b, _ := util.JsonEncode(cfg)
		Instance.items[pointer].m = cfg
		Instance.items[pointer].hash = util.Md5String(string(b))
		Instance.items[pointer].elem = reflect.ValueOf(cfg).Elem()
	}

	if v := Instance.items[pointer]; v != nil {
		reflect.ValueOf(cfg).Elem().Set(v.elem)
	}
	return cfg
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
