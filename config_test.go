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
)

const TEST_PARSER = "mem" // mem、toml、 toml_apollo、toml_nacos

type RedisCfg struct {
	DialAddr string `toml:"addr"`
	DialPawd string `toml:"pawd"`
}

type RedisConfig struct {
	Cfgs map[string]*RedisCfg `toml:"redis" antcfg:"redis"`
}

type MysqlCfg struct {
	DialUser string `toml:"user"`
	DialPawd string `toml:"pawd"`
	DialHost string `toml:"host"`
	DialPort string `toml:"port"`
	DialName string `toml:"name"`
}

type MysqlConfig struct {
	Cfgs map[string]*MysqlCfg `toml:"mysql" antcfg:"mysql"`
}

func TestMain(m *testing.M) {
	var p parser.Parser
	var source string
	var mem interface{}

	switch TEST_PARSER {
	case "mem":
		p = parser.NewMemParser()
		mem = &struct {
			Name  string
			Redis map[string]*RedisCfg `antcfg:"redis"`
			Mysql map[string]*MysqlCfg `antcfg:"mysql"`
		}{
			Name: "app1",
			Redis: map[string]*RedisCfg{
				"stats": {DialAddr: "rd1addr", DialPawd: "rd1pawd"},
				"uaap":  {DialAddr: "rd2addr", DialPawd: "rd2pawd"},
			},
			Mysql: map[string]*MysqlCfg{
				"default": {DialUser: "md1addr", DialPawd: "md1pawd"},
				"md2":     {DialUser: "md2addr", DialPawd: "md2pawd"},
			},
		}
	case "toml":
		p = parser.NewTomlParser()
		source = "./test/toml_test.toml"
	case "toml_apollo":
		p = parser.NewTomlApolloParser()
		source = "./test/toml_apollo_test.toml"
	case "toml_nacos":
		p = parser.NewTomlNacosParser()
		source = "./test/toml_nacos_test.toml"
	}

	NewConfig(p,
		options.WithCfgSource(source),
		options.WithMemoryVariable(mem),
		options.WithCheckInterval(10),
		options.WithOnChangeFn(func(cfg interface{}) {
			switch v := cfg.(type) {
			case *RedisConfig:
				fmt.Println("global change redis: ", v.Cfgs["stats"], v.Cfgs["uaap"])
			case *MysqlConfig:
				fmt.Println("global change mysql: ", v.Cfgs["default"])
			default:
				fmt.Println("global change: ", v)
			}
		}),
		options.WithOnErrorFn(func(e error) {
			fmt.Println("error: ", e)
		}))

	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	redisCfg := &RedisConfig{}
	mysqlCfg := &MysqlConfig{}

	fmt.Println("redis: ", Get(redisCfg, options.WithOpOnChangeFn(func(redisCfg interface{}) {
		fmt.Println("private change redis : ", redisCfg)
	})).(*RedisConfig).Cfgs["stats"])

	fmt.Println("redis: ", Get(redisCfg, options.WithOpOnChangeFn(func(redisCfg interface{}) {
		fmt.Println("private change redis : ", redisCfg)
	})).(*RedisConfig).Cfgs["uaap"])

	fmt.Println("mysql: ", Get(mysqlCfg).(*MysqlConfig).Cfgs["default"])

	time.Sleep(1 * time.Hour)
}
