# Config

支持扩展的配置文件解析库

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ant-libs-go/config?status.png)](http://godoc.org/github.com/ant-libs-go/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/ant-libs-go/config)](https://goreportcard.com/report/github.com/ant-libs-go/config)

# 特性

* 针对 toml 格式配置文件，支持 import 关键字，可以引入其他 toml 文件
* 支持定时检测配置文件更新状态，完成 reload 之后发起回调通知
* 支持配置文件加载失败时的回调通知
* 当前支持本地 toml 格式、Apollo toml、Nacos toml 配置方案，如需ini、yaml等配置格式可按协议进行实现

## 安装

	go get github.com/ant-libs-go/config

# 快速开始

* toml 配置
 
    ```
    import = ["app-local"] // 同目录下的 app-loca.toml，按顺序读取

    [mysql.default]
        user = "root"
        pawd = "123456"
        host = "127.0.0.1"
        port = "3306"
        name = "business"
    [redis.default]
        addr = "127.0.0.1:6379"
        pawd = "123456"
    ```

* 使用方法

    ```golang
    type MysqlConfig struct {
        Cfgs map[string]*struct {
            DialUser string `toml:"user"`
            DialPawd string `toml:"pawd"`
            DialHost string `toml:"host"`
            DialPort string `toml:"port"`
            DialName string `toml:"name"`
        } `toml:"mysql"`
    }

    func main() {
        config.New(parser.NewTomlParser(),
            options.WithCfgSource("./test/toml_test.toml"),
            options.WithCheckInterval(10),
            options.WithOnChangeFn(func(data interface{}) { // 配置发生变化时触发
                fmt.Println("change.....")
                switch v := cfg.(type) {
                case *RedisConfig:
                    fmt.Println("change redis: ", v.Cfgs["default"])
                case *MysqlConfig:
                    fmt.Println("change mysql: ", v.Cfgs["default"])
                }
            }),
            options.WithOnErrorFn(func(err error) { // 加载配置出现错误时触发
                fmt.Println("err", err)
            }))

        fmt.Printf("ret: %+v\n", config.Get(&RedisConfig{}).(*RedisConfig))
        fmt.Printf("ret: %+v\n", config.Get(&RedisConfig{}).(*RedisConfig).Cfgs["default"]))
        
        time.Sleep(1 * time.Hour)
    }
    ```
