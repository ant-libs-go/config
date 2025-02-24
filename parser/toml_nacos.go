/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2021-10-12 15:58:25
# File Name: toml_nacos.go
# Description:
####################################################################### */

package parser

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/util"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type TomlNacosConfig struct {
	Addr        string `toml:"addr"`
	DataId      string `toml:"data_id"`
	GroupId     string `toml:"group_id"`
	NamespaceId string `toml:"namespace_id"`
	AccessKey   string `toml:"access_key"`
	SecretKey   string `toml:"secret_key"`
	CacheDir    string `toml:"cache_dir"`
	LogDir      string `toml:"log_dir"`
	LogLevel    string `toml:"log_level"` // debug, info, warn, error
}

type TomlNacosParser struct {
	initOnce   sync.Once
	listenOnce sync.Once

	cli      config_client.IConfigClient
	modTime  int64
	entrance *TomlNacosConfig
}

func NewTomlNacosParser() *TomlNacosParser {
	o := &TomlNacosParser{}
	return o
}

func (this *TomlNacosParser) Unmarshal(cfg interface{}, opts *options.Options) (err error) {
	this.initOnce.Do(func() {
		var host, port string
		this.entrance = &TomlNacosConfig{}
		if len(opts.Sources) == 0 {
			err = fmt.Errorf("config source not specified")
		}
		if err == nil {
			_, err = toml.DecodeFile(opts.Sources[0], this.entrance)
		}
		if err == nil {
			host, port, err = net.SplitHostPort(this.entrance.Addr)
		}
		if err == nil {
			pwd, _ := os.Getwd()
			this.cli, err = clients.NewConfigClient(vo.NacosClientParam{
				ServerConfigs: []constant.ServerConfig{
					*constant.NewServerConfig(host, uint64(util.StrToInt64(port, 80))),
				},
				ClientConfig: constant.NewClientConfig(
					constant.WithTimeoutMs(5000),
					constant.WithNamespaceId(this.entrance.NamespaceId),
					constant.WithNotLoadCacheAtStart(true),
					constant.WithCacheDir(util.AbsPath(this.entrance.CacheDir, pwd)),
					constant.WithLogDir(util.AbsPath(this.entrance.LogDir, pwd)),
					constant.WithLogLevel(this.entrance.LogLevel),
					constant.WithAccessKey(this.entrance.AccessKey),
					constant.WithSecretKey(this.entrance.SecretKey),
				),
			})
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
		for _, source := range sources {
			if IsLocalFile(source) == true {
				continue
			}
			err = this.cli.ListenConfig(vo.ConfigParam{
				DataId:   source,
				Group:    this.entrance.GroupId,
				OnChange: func(namespace, group, dataId, data string) { this.modTime = time.Now().Unix() },
			})
		}
	})

	for _, source := range sources {
		if err = this.decode(cfg, source); err != nil {
			return
		}
	}
	return
}

func (this *TomlNacosParser) GetLastModTime(opts *options.Options) (r int64, err error) {
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

func (this *TomlNacosParser) parseSource(opts *options.Options) (r []string, err error) {
	r = []string{}
	r = append(r, opts.Sources...)

	var str string
	t := &TomlImport{}

	str, err = this.cli.GetConfig(vo.ConfigParam{DataId: this.entrance.DataId, Group: this.entrance.GroupId})
	if err == nil {
		_, err = toml.Decode(str, t)
	}
	if err != nil {
		return
	}

	for _, v := range t.Import {
		r = append(r, strings.TrimSpace(v))
	}
	r = append(r, this.entrance.DataId)
	return
}

func (this *TomlNacosParser) decode(cfg interface{}, source string) (err error) {
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
		if str, err = this.cli.GetConfig(vo.ConfigParam{DataId: source, Group: this.entrance.GroupId}); err != nil {
			err = fmt.Errorf("nacos get config fail, %s", err)
		}
	}
	if err == nil {
		if _, err = toml.Decode(str, cfg); err != nil {
			err = fmt.Errorf("nacos config source decode fail, %s", err)
		}
	}
	return
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
