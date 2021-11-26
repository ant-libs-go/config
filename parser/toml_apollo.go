/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2021-09-09 17:33:02
# File Name: toml_apollo.go
# Description:
####################################################################### */

package parser

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ant-libs-go/config/options"
	"github.com/philchia/agollo/v4"
)

type TomlApolloConfig struct {
	AppId    string `toml:"app_id"`
	CacheDir string `toml:"cache_dir"`
	MetaAddr string `toml:"meta_addr"`
}

type TomlApolloParser struct {
	initOnce   sync.Once
	listenOnce sync.Once

	modTime  int64
	entrance *TomlApolloConfig
}

func NewTomlApolloParser() *TomlApolloParser {
	o := &TomlApolloParser{}
	return o
}

func (this *TomlApolloParser) Unmarshal(cfg interface{}, opts *options.Options) (err error) {
	this.initOnce.Do(func() {
		this.entrance = &TomlApolloConfig{}
		if len(opts.Sources) == 0 {
			err = fmt.Errorf("config source not specified")
		}
		if err == nil {
			_, err = toml.DecodeFile(opts.Sources[0], this.entrance)
		}
		if err == nil {
			err = agollo.Start(&agollo.Conf{
				AppID:    this.entrance.AppId,
				CacheDir: this.entrance.CacheDir,
				MetaAddr: this.entrance.MetaAddr})
		}
	})

	if err != nil {
		return
	}

	var sources []string
	if sources, err = this.parseSource(opts); err != nil {
		return
	}

	this.listenOnce.Do(func() {
		agollo.OnUpdate(func(e *agollo.ChangeEvent) { this.modTime = time.Now().Unix() })

		for _, source := range sources {
			if IsLocalFile(source) == true {
				continue
			}
			agollo.SubscribeToNamespaces(strings.TrimSpace(source))
		}
	})

	for _, source := range sources {
		if err = this.decode(cfg, source); err != nil {
			return
		}
	}
	return
}

func (this *TomlApolloParser) GetLastModTime(opts *options.Options) (r int64, err error) {
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

func (this *TomlApolloParser) parseSource(opts *options.Options) (r []string, err error) {
	r = []string{}
	r = append(r, opts.Sources...)

	for _, source := range strings.Split(agollo.GetString("import"), ",") {
		r = append(r, strings.TrimSpace(source))
	}
	return
}

func (this *TomlApolloParser) decode(cfg interface{}, source string) (err error) {
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

	var str string
	if err == nil {
		str = agollo.GetContent(agollo.WithNamespace(source))
	}
	if err == nil {
		if _, err = toml.Decode(str, cfg); err != nil {
			err = fmt.Errorf("apollo config source decode fail, %s", err)
		}
	}
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
