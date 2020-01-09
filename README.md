# config
一个读取配置的包

# 功能
 - 支持 import 关键字，可以引入其他 toml 文件。再也不用将所有配置文件写在一个文件中
 - 支持定时检测配置文件更新状态，完成 reload 之后发起回调通知
 - 支持配置文件加载失败时的回调通知
 - 目前仅支持本地 toml 格式配置文件

# 基本使用
 - toml 配置
 
 	```
	debug = true
	port = ":8081"
	logFile = "conf/log.xml"
	import = ["app-local"] // 同目录下的 app-loca.toml，按顺序读取

	[gc]
		percent = 100
	[db]
		[db.master]
			host = "123"
			port = "123"
			user = "123"
			pawd = "123"
			name = "123"
		[db.slave]
			host = "123"
			port = "123"
			user = "123"
			pawd = "123"
			name = "123"

- 使用方法

	```golang
    type AppConfig struct {
        Debug   bool
        Port    string
        LogFile string

        Gc *struct {
            Percent int // default 100
        }

        Db map[string]*struct {
            Host string
            Port string
            User string
            Pawd string
            Name string
        }
    }

	func main() {
		Default, _ := config.New(&AppConfig{}, parser.NewTomlParser(),
			options.WithCfgFile(*cfg),
			options.WithCheckInterval(1),
			options.WithOnChangeFn(func(data interface{}) { // 配置发生变化时触发
				fmt.Println("change")
				fmt.Printf("ret: %+v\n", data.(*AppConfig))
			}),
			options.WithOnErrorFn(func(err error) { fmt.Println("err", err) }), // 加载配置出现错误时触发
		)
		
		fmt.Printf("ret: %+v\n", Default.Get().(*AppConfig))
		fmt.Printf("ret: %+v\n", Default.Get().(*AppConfig).Db["slave"]))
		
		time.Sleep(1000 * time.Second)
    }
```

# 进一步使用
 - 可以参考 parser interface，自定义其他类型的配置，如 ini、apollo（携程的开源产品）
