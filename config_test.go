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

const TEST_PARSER = "toml_apollo" // toml / toml_apollo

func TestMain(m *testing.M) {
	switch TEST_PARSER {
	case "toml":
		NewConfig(parser.NewTomlParser(),
			options.WithCfgSource("./test/toml_test.toml"),
			options.WithCheckInterval(10),
			options.WithOnChangeFn(func(cfg interface{}) {
				fmt.Println("change.....")
				switch v := cfg.(type) {
				case *RedisConfig:
					fmt.Println("change redis: ", v.Cfgs["default"])
				case *MysqlConfig:
					fmt.Println("change mysql: ", v.Cfgs["default"])
				}
			}),
			options.WithOnErrorFn(func(e error) {
				fmt.Println("error: ", e)
			}))
	case "toml_apollo":
		NewConfig(parser.NewTomlApolloParser(),
			options.WithCfgSource("./test/toml_apollo_test.toml"),
			options.WithCheckInterval(10),
			options.WithOnChangeFn(func(cfg interface{}) {
				fmt.Println("change.....")
				switch v := cfg.(type) {
				case *RedisConfig:
					fmt.Println("change redis: ", v.Cfgs["default"])
				case *MysqlConfig:
					fmt.Println("change mysql: ", v.Cfgs["default"])
				}
			}),
			options.WithOnErrorFn(func(e error) {
				fmt.Println("error: ", e)
			}))
	}
	os.Exit(m.Run())
}

type RedisConfig struct {
	Cfgs map[string]*struct {
		DialAddr     string `toml:"addr"`
		DialUsername string `toml:"user"`
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

func TestBasic(t *testing.T) {
	cfg := &RedisConfig{}
	fmt.Println("redis: ", Get(cfg).(*RedisConfig).Cfgs["default"])
	fmt.Println("mysql: ", Get(&MysqlConfig{}).(*MysqlConfig).Cfgs["default"])
	time.Sleep(1 * time.Hour)
}
