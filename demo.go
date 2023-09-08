package main

import (
	"flag"
	"fmt"

	"demo/demo"
	"demo/internal/config"
	"demo/internal/server"
	"demo/internal/svc"

	"context"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net/http"
)

var configFile = flag.String("f", "etc/demo.yaml", "the config file")

func limit(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var n = 100

	l := syncx.NewLimit(n)

	if l.TryBorrow() {
		defer func() {
			if err := l.Return(); err != nil {
				logx.Error(err)
			}
		}()
		return handler(ctx, req)
	} else {
		logx.Errorf("concurrent connections over %d, rejected with code %d",
			n, http.StatusServiceUnavailable)
		return nil, status.Error(codes.Unavailable, "concurrent connections over limit")
	}
}
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		demo.RegisterDemoServer(grpcServer, server.NewDemoServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(limit)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
