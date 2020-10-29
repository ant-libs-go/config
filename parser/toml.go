/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:59:22
# File Name: toml.go
# Description:
####################################################################### */

package parser

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/ant-libs-go/config/options"
)

type TomlParser struct{}

func NewTomlParser() *TomlParser {
	return &TomlParser{}
}

type Import struct {
	Import []string
}

func (this *TomlParser) parseSource(opts *options.Options) (r []string, err error) {
	r = []string{}
	r = append(r, opts.Sources...)

	for _, source := range r {
		imp := &Import{}
		if err = this.load(imp, source); err != nil {
			return
		}
		for _, v := range imp.Import {
			r = append(r, path.Join(path.Dir(source), v+".toml"))
		}
	}
	return
}

func (this *TomlParser) Unmarshal(cfg interface{}, opts *options.Options) (err error) {
	var sources []string
	if sources, err = this.parseSource(opts); err != nil {
		return
	}
	for _, source := range sources {
		if err = this.load(cfg, source); err != nil {
			return
		}
	}
	return
}

func (this *TomlParser) load(cfg interface{}, source string) (err error) {
	if len(source) == 0 {
		err = fmt.Errorf("config source not specified")
		return
	}
	_, err = toml.DecodeFile(source, cfg)
	if err != nil {
		err = fmt.Errorf("config source unmarshal fail, %s", err)
		return
	}
	return
}

func (this *TomlParser) GetLastModTime(opts *options.Options) (r int64, err error) {
	var sources []string
	if sources, err = this.parseSource(opts); err != nil {
		return
	}
	for _, source := range sources {
		var modTime int64
		if modTime, err = this.getLastModTime(source); err != nil {
			return
		}
		if modTime > r {
			r = modTime
		}
	}
	return
}

func (this *TomlParser) getLastModTime(source string) (r int64, err error) {
	fd, err := os.Stat(source)
	if err != nil {
		err = fmt.Errorf("get source last modified time fail, %s", err)
		return
	}
	r = fd.ModTime().Unix()
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
