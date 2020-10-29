/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-29 19:33:22
# File Name: config_mgr.go
# Description:
####################################################################### */

package config

import (
	"fmt"

	"github.com/ant-libs-go/config/options"
)

var Default *Config

func Load(cfg interface{}, opts ...options.Option) (r *Config, err error) {
	if Default == nil {
		err = fmt.Errorf("instance was not initialized")
		return
	}
	return Default.Load(cfg, opts...)
}

func Get(cfg interface{}) interface{} {
	if Default == nil {
		return nil
	}
	return Default.Get(cfg)
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
