/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-01-05 12:59:22
# File Name: toml.go
# Description:
####################################################################### */

package parser

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/ant-libs-go/config/options"
)

type TomlParser struct {
	modTime int64
}

func NewTomlParser() *TomlParser {
	o := &TomlParser{}
	return o
}

func (this *TomlParser) Unmarshal(cfg interface{}, opts *options.Options) (err error) {
	var sources []string
	if sources, err = this.parseSource(opts); err != nil {
		return
	}

	for _, source := range sources {
		if err = this.decode(cfg, source); err != nil {
			return
		}
	}
	return
}

func (this *TomlParser) GetLastModTime(opts *options.Options) (r int64, err error) {
	var sources []string
	if sources, err = this.parseSource(opts); err != nil {
		return
	}

	for _, source := range sources {
		if IsLocalFile(source) == false {
			continue
		}
		var modTime int64
		if modTime, err = ParseFileLastModTime(source); err != nil {
			return
		}
		if modTime > this.modTime {
			this.modTime = modTime
		}
	}
	return this.modTime, nil
}

func (this *TomlParser) parseSource(opts *options.Options) (r []string, err error) {
	r = []string{}
	r = append(r, opts.Sources...)

	t := &TomlImport{}
	dir, _ := path.Split(opts.Sources[0])

	if err = this.decode(t, opts.Sources[0]); err != nil {
		return
	}

	for _, v := range t.Import {
		r = append(r, fmt.Sprintf("%s%s.toml", dir, v))
	}
	return
}

func (this *TomlParser) decode(cfg interface{}, source string) (err error) {
	if len(source) == 0 {
		err = fmt.Errorf("config source not specified")
		return
	}

	if IsLocalFile(source) == true {
		if _, err = toml.DecodeFile(source, cfg); err != nil {
			err = fmt.Errorf("local config source decode fail, %s", err)
		}
		return
	}

	err = fmt.Errorf("local config source[%s] not found", source)
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
