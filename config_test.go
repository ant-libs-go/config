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

const TEST_PARSER = "toml_nacos" // toml、 toml_apollo、toml_nacos

type RedisConfig struct {
	Cfgs map[string]*struct {
		DialAddr string `toml:"addr"`
		DialPawd string `toml:"pawd"`
	} `toml:"redis"`
}

type MysqlConfig struct {
	Cfgs map[string]*struct {
		DialUser string `toml:"user"`
		DialPawd string `toml:"pawd"`
		DialHost string `toml:"host"`
		DialPort string `toml:"port"`
		DialName string `toml:"name"`
	} `toml:"mysql"`
}

func TestMain(m *testing.M) {
	var p parser.Parser
	var source string

	switch TEST_PARSER {
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
