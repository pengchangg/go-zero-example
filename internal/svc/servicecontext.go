package svc

import (
	"demo/internal/config"
	"demo/model"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	SqlConn     sqlx.SqlConn
	RedisClient *redis.Redis
	UserModel   model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	SqlConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:      c,
		SqlConn:     SqlConn,
		RedisClient: redis.MustNewRedis(c.Redis.RedisConf),
		UserModel:   model.NewUsersModel(SqlConn, c.Cache),
	}
}
