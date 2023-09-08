package main

import (
	"context"
	"demo/demo"
	"demo/internal/config"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"log"
	"time"
)

var configFile = flag.String("f", "etc/demo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	conn := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: c.Etcd,
	})

	client := demo.NewDemoClient(conn.Conn())
	for {
		resp, err := client.Ping(
			context.Background(),
			&demo.Request{
				Ping: fmt.Sprintf("ping.time = %s", time.Now().Format(time.DateTime)),
			},
		)
		if err != nil {
			log.Printf("err--%s \n", err.Error())
		} else {
			log.Println(resp)
		}
		time.Sleep(time.Second * 2)
	}
}
