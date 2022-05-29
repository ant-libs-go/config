/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2022-05-29 20:56:25
# File Name: parser/mem.go
# Description:
####################################################################### */

package parser

import (
	"reflect"

	"github.com/ant-libs-go/config/options"
)

type MemParser struct{}

func NewMemParser() *MemParser {
	o := &MemParser{}
	return o
}

func (this *MemParser) Unmarshal(cfg interface{}, opts *options.Options) (err error) {
	m := map[string]reflect.Value{}
	valMem := reflect.ValueOf(opts.MemoryVariable).Elem()
	for i := 0; i < valMem.NumField(); i++ {
		tag := valMem.Type().Field(i).Tag.Get("antcfg")
		if len(tag) == 0 {
			continue
		}
		m[tag] = reflect.ValueOf(valMem.Field(i).Interface())
	}

	valCfg := reflect.ValueOf(cfg).Elem()
	for i := 0; i < valCfg.NumField(); i++ {
		tag := valCfg.Type().Field(i).Tag.Get("antcfg")
		if len(tag) == 0 {
			continue
		}
		if _, ok := m[tag]; !ok {
			continue
		}
		valCfg.Field(i).Set(m[tag])
	}
	return
}

func (this *MemParser) GetLastModTime(opts *options.Options) (r int64, err error) {
	return 0, nil
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
