package app

import (
	"context"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"

	"github.com/gin-gonic/gin"
)

// SQLGetTAppConfigIntByK 查询配置
func SQLGetTAppConfigIntByK(ctx context.Context, tx hcommon.DbExeAble, k string) (*model.DBTAppConfigInt, error) {
	var row model.DBTAppConfigInt
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    k,
    v
FROM
	t_app_config_int
WHERE
	k=:k
LIMIT 1`,
		gin.H{
			"k": k,
		},
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLGetTAppStatusIntByK 查询配置
func SQLGetTAppStatusIntByK(ctx context.Context, tx hcommon.DbExeAble, k string) (*model.DBTAppStatusInt, error) {
	var row model.DBTAppStatusInt
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    k,
    v
FROM
	t_app_status_int
WHERE
	k=:k
LIMIT 1`,
		gin.H{
			"k": k,
		},
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLGetTAddressKeyFreeCount 获取剩余可用地址数
func SQLGetTAddressKeyFreeCount(ctx context.Context, tx hcommon.DbExeAble) (int64, error) {
	var count int64
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&count,
		`SELECT
	IFNULL(COUNT(*), 0)
FROM
	t_address_key
WHERE
	use_tag=0`,
		gin.H{},
	)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return count, nil
}
