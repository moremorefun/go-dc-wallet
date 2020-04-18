package app

import (
	"context"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"strings"

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

// SQLSelectTAddressKeyColByAddress 根据ids获取
func SQLSelectTAddressKeyColByAddress(ctx context.Context, tx hcommon.DbExeAble, cols []string, addresses []string) ([]*model.DBTAddressKey, error) {
	if len(addresses) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
WHERE
	address IN (:addresses)
	AND use_tag>0`)

	var rows []*model.DBTAddressKey
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"addresses": addresses,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppStatusIntByK 更新
func SQLUpdateTAppStatusIntByK(ctx context.Context, tx hcommon.DbExeAble, row *model.DBTAppStatusInt) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_status_int
SET
    v=:v
WHERE
	k=:k`,
		gin.H{
			"k": row.K,
			"v": row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}
