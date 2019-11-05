package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

type Data struct {
	domain string
	value  string
}

func (data *Data) test(c chan bool) string {

	c <- false
	if 1 == 1 {
		return "sdfs"
	}
	fmt.Print(13213)
	return "sdf"

}

func main() {

	var domain, value, beginIp, endIp string
	var port int64
	var thread int
	app := cli.NewApp()

	app.Name = "FUCK CDN"                                         // 指定程序名称
	app.Usage = "fuckcdn --domain baidu.com --port 80 --vaule 百度" //  程序功能描述
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "domain,d", // 配置名称
			Usage:       "域名",       // 配置描述
			Destination: &domain,
		},
		cli.Int64Flag{
			Name:        "port, p", // 配置名称
			Value:       80,        // 缺省配置值
			Usage:       "端口",      // 配置描述
			Destination: &port,     // 保存配置值
		},
		cli.StringFlag{
			Name:        "value", // 配置名称
			Usage:       "关键字",   // 配置描述
			Destination: &value,  // 保存配置值
		},
		cli.IntFlag{
			Name:        "thread, t", // 配置名称
			Value:       200,         // 缺省配置值
			Usage:       "端口",        // 配置描述
			Destination: &thread,     // 保存配置值
		},
		cli.StringFlag{
			Name:        "beginip, b", // 配置名称
			Value:       "1.1.1.1",    // 缺省配置值
			Usage:       "开始IP",       // 配置描述
			Destination: &beginIp,     // 保存配置值
		},
		cli.StringFlag{
			Name:        "endip, e",        // 配置名称
			Value:       "255.255.255.255", // 缺省配置值
			Usage:       "结束IP",            // 配置描述
			Destination: &endIp,            // 保存配置值
		},
	}
	app.Action = func(c *cli.Context) error {
		if domain != "" && value != "" {
			begin := NormalIpToten(beginIp)
			end := NormalIpToten(endIp)
			if begin > end {
				temp := begin
				begin = end
				end = temp
			}
			fmt.Print(begin, end)
			data := &Data{domain, value}
			c := make(chan string, 100000)
			go data.ReciveMessage(c)
			data.Start(port, begin, c, thread, end)
		}
		return nil
	}

	app.Run(os.Args)

}
