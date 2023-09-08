package logic

import (
	"context"

	"demo/demo"
	"demo/internal/svc"

	"demo/model"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/exp/rand"
	"time"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(req *demo.Request) (*demo.Response, error) {

	rand.Seed(uint64(time.Now().UnixNano()))
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, 10)

	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	lock := redis.NewRedisLock(l.svcCtx.RedisClient, "test-lock")
	lock.SetExpire(10)
	acquire, err := lock.AcquireCtx(l.ctx)
	switch {
	case err != nil:
		return nil, err
	case acquire:
		defer lock.Release() // 释放锁
		_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.Users{
			Age:  rand.Int63n(81),
			Name: string(randomString),
		})

		if err != nil {
			return nil, err
		}

		limit100, err := l.svcCtx.UserModel.ListLimit100(l.ctx)
		if err != nil {
			logx.Errorf("ListLimit100 error: %+v", err)
			return nil, err
		}

		return &demo.Response{
			Pong: fmt.Sprintf("ping.req.msg =%s , limit100 = %+v", req.Ping, limit100),
		}, nil
	case !acquire:
		return nil, errors.New("没有获取到lock")
	}
	return nil, errors.New("empty")
}
