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

// SQLGetTAppConfigStrByK 查询配置
func SQLGetTAppConfigStrByK(ctx context.Context, tx hcommon.DbExeAble, k string) (*model.DBTAppConfigStr, error) {
	var row model.DBTAppConfigStr
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    k,
    v
FROM
	t_app_config_str
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

// SQLSelectTTxColByOrg 获取未整理交易
func SQLSelectTTxColByOrg(ctx context.Context, tx hcommon.DbExeAble, cols []string) ([]*model.DBTTx, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
WHERE
	org_status=:org_status`)

	var rows []*model.DBTTx
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"org_status": 0,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLGetTSendMaxNonce 获取地址的nonce
func SQLGetTSendMaxNonce(ctx context.Context, tx hcommon.DbExeAble, address string) (int64, error) {
	var i int64
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&i,
		`SELECT 
	IFNULL(MAX(nonce), -1)
FROM
	t_send
WHERE
	from_address=:address
LIMIT 1`,
		gin.H{
			"address": address,
		},
	)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}
	return i + 1, nil
}

// SQLGetTAddressKeyColByAddress 根据address查询
func SQLGetTAddressKeyColByAddress(ctx context.Context, tx hcommon.DbExeAble, cols []string, address string) (*model.DBTAddressKey, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
WHERE
	address=:address`)

	var row model.DBTAddressKey
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"address": address,
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

// SQLUpdateTTxOrgStatusByAddresses 更新
func SQLUpdateTTxOrgStatusByAddresses(ctx context.Context, tx hcommon.DbExeAble, addresses []string, row model.DBTTx) (int64, error) {
	if len(addresses) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx
SET
    org_status=:org_status,
    org_msg=:org_msg,
    org_time=:org_time
WHERE
	to_address IN (:addresses)`,
		gin.H{
			"addresses":  addresses,
			"org_status": row.OrgStatus,
			"org_msg":    row.OrgMsg,
			"org_time":   row.OrgTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLUpdateTSendStatusByTxIDs 更新
func SQLUpdateTSendStatusByTxIDs(ctx context.Context, tx hcommon.DbExeAble, txIDs []string, row model.DBTSend) (int64, error) {
	if len(txIDs) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_send
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	tx_id IN (:tx_ids)`,
		gin.H{
			"tx_ids":        txIDs,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLUpdateTSendStatusByIDs 更新
func SQLUpdateTSendStatusByIDs(ctx context.Context, tx hcommon.DbExeAble, ids []int64, row model.DBTSend) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_send
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id IN (:ids)`,
		gin.H{
			"ids":           ids,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLSelectTSendColByStatus 根据ids获取
func SQLSelectTSendColByStatus(ctx context.Context, tx hcommon.DbExeAble, cols []string, status int64) ([]*model.DBTSend, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send
WHERE
	handle_status=:handle_status`)

	var rows []*model.DBTSend
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"handle_status": status,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
