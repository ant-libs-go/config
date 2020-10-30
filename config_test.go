/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2020-10-29 13:50:17
# File Name: config_test.go
# Description:
####################################################################### */

package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ant-libs-go/config/options"
	"github.com/ant-libs-go/config/parser"
	. "github.com/smartystreets/goconvey/convey"
)

var globalCfg *Config

func TestMain(m *testing.M) {
	globalCfg = New(parser.NewTomlParser(),
		options.WithCfgSource("/tmp/app.toml"),
		options.WithCheckInterval(1),
		options.WithOnChangeFn(func(cfg interface{}) { fmt.Println(cfg.(*RedisConfig).Redis) }),
		options.WithOnErrorFn(func(e error) { fmt.Println(e) }))
	os.Exit(m.Run())
}

type RedisConfig struct {
	Redis *struct {
		Addrs []string `toml:"addrs"`
		Pawd  string   `toml:"pawd"`
	} `toml:"redis"`
}

type MysqlConfig struct {
	Mysql map[string]*struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
		User string `toml:"user"`
		Pawd string `toml:"pawd"`
		Name string `toml:"name"`
	} `toml:"mysql"`
}

func TestBasic(t *testing.T) {
	_, err := globalCfg.Load(&RedisConfig{}, options.WithOnChangeFn(func(cfg interface{}) { fmt.Println("redis===>", cfg.(*RedisConfig).Redis) }))
	fmt.Println(err)
	_, err = globalCfg.Load(&MysqlConfig{}, options.WithOnChangeFn(func(cfg interface{}) { fmt.Println("mysql===>", cfg.(*MysqlConfig).Mysql) }))
	fmt.Println(err)
	fmt.Println(globalCfg.Get(&RedisConfig{}).(*RedisConfig).Redis)
	cfg := &MysqlConfig{}
	fmt.Println(globalCfg.Get(cfg).(*MysqlConfig).Mysql)
	fmt.Println("-->", cfg, cfg.Mysql)
	time.Sleep(100 * time.Second)

	Convey("TestBasic", t, func() {
		Convey("load should return nil when config exist", func() {
			So(err, ShouldBeNil)
		})
		Convey("the value should be equal to the defined string", func() {
			So(globalCfg.Get(&RedisConfig{}).(*RedisConfig).Redis.Pawd, ShouldEqual, "ddddd")
		})
	})
}
