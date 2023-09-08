package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ EntTypesModel = (*customEntTypesModel)(nil)

type (
	// EntTypesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEntTypesModel.
	EntTypesModel interface {
		entTypesModel
	}

	customEntTypesModel struct {
		*defaultEntTypesModel
	}
)

// NewEntTypesModel returns a model for the database table.
func NewEntTypesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) EntTypesModel {
	return &customEntTypesModel{
		defaultEntTypesModel: newEntTypesModel(conn, c, opts...),
	}
}
