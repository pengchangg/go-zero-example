package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		ListLimit100(ctx context.Context) ([]Users, error)
		Page(ctx context.Context, page int) ([]Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}

func (c *defaultUsersModel) ListLimit100(ctx context.Context) ([]Users, error) {
	query := fmt.Sprintf("select %s from %s order by id desc limit 100 ", usersRows, c.table)
	var resp []Users
	err := c.QueryRowsNoCacheCtx(ctx, &resp, query)
	return resp, err
}

func (c *defaultUsersModel) Page(ctx context.Context, page int) ([]Users, error) {
	return nil, nil
}
