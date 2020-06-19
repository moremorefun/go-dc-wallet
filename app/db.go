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
func SQLGetTAddressKeyFreeCount(ctx context.Context, tx hcommon.DbExeAble, symbol string) (int64, error) {
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
	use_tag=0
	AND symbol=:symbol`,
		gin.H{
			"symbol": symbol,
		},
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
	address IN (:addresses)`)

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
func SQLSelectTTxColByOrg(ctx context.Context, tx hcommon.DbExeAble, cols []string, orgStatus int64) ([]*model.DBTTx, error) {
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
			"org_status": orgStatus,
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

// SQLGetTSendPendingBalance 获取地址的打包数额
func SQLGetTSendPendingBalance(ctx context.Context, tx hcommon.DbExeAble, address string) (int64, error) {
	var i int64
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&i,
		`SELECT 
	IFNULL(SUM(balance), 0)
FROM
	t_send
WHERE
	from_address=:address
	AND handle_status<2
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
	return i, nil
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

// SQLUpdateTTxOrgStatusByIDs 更新
func SQLUpdateTTxOrgStatusByIDs(ctx context.Context, tx hcommon.DbExeAble, ids []int64, row model.DBTTx) (int64, error) {
	if len(ids) == 0 {
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
	id IN (:ids)`,
		gin.H{
			"ids":        ids,
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

// SQLUpdateTTxStatusByIDs 更新
func SQLUpdateTTxStatusByIDs(ctx context.Context, tx hcommon.DbExeAble, ids []int64, row model.DBTTx) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx
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

// SQLUpdateTTxErc20StatusByIDs 更新
func SQLUpdateTTxErc20StatusByIDs(ctx context.Context, tx hcommon.DbExeAble, ids []int64, row model.DBTTxErc20) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_erc20
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

// SQLUpdateTWithdrawStatusByTxIDs 更新
func SQLUpdateTWithdrawStatusByTxIDs(ctx context.Context, tx hcommon.DbExeAble, txIDs []string, row model.DBTWithdraw) (int64, error) {
	if len(txIDs) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	tx_hash IN (:tx_ids)`,
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

// SQLSelectTWithdrawColByStatus 根据ids获取
func SQLSelectTWithdrawColByStatus(ctx context.Context, tx hcommon.DbExeAble, cols []string, status int64, symbols []string) ([]*model.DBTWithdraw, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	handle_status=:handle_status
	AND symbol IN (:symbols)`)

	var rows []*model.DBTWithdraw
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"handle_status": status,
			"symbols":       symbols,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLGetTWithdrawColForUpdate 根据id查询
func SQLGetTWithdrawColForUpdate(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64, status int64) (*model.DBTWithdraw, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	id=:id
	AND handle_status=:handle_status
FOR UPDATE`)

	var row model.DBTWithdraw
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"id":            id,
			"handle_status": status,
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

// SQLUpdateTWithdrawGenTx 更新
func SQLUpdateTWithdrawGenTx(ctx context.Context, tx hcommon.DbExeAble, row *model.DBTWithdraw) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
SET
    tx_hash=:tx_hash,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"tx_hash":       row.TxHash,
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

// SQLUpdateTWithdrawStatusByIDs 更新
func SQLUpdateTWithdrawStatusByIDs(ctx context.Context, tx hcommon.DbExeAble, ids []int64, row *model.DBTWithdraw) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
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

// SQLGetTAppLockColByK 根据id查询
func SQLGetTAppLockColByK(ctx context.Context, tx hcommon.DbExeAble, cols []string, k string) (*model.DBTAppLock, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
WHERE
	k=:k
	AND v=1`)

	var row model.DBTAppLock
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
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

// SQLCreateTAppLockUpdate 创建
func SQLCreateTAppLockUpdate(ctx context.Context, tx hcommon.DbExeAble, row *model.DBTAppLock) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_lock (
    id,
    k,
    v,
    create_time
) VALUES (
    :id,
    :k,
    :v,
    :create_time
) ON DUPLICATE KEY UPDATE 
	v=:v,
	create_time=:create_time`,
			gin.H{
				"id":          row.ID,
				"k":           row.K,
				"v":           row.V,
				"create_time": row.CreateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_lock (
    k,
    v,
    create_time
) VALUES (
    :k,
    :v,
    :create_time
) ON DUPLICATE KEY UPDATE 
	v=:v,
	create_time=:create_time`,
			gin.H{
				"k":           row.K,
				"v":           row.V,
				"create_time": row.CreateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLUpdateTAppLockByK 更新
func SQLUpdateTAppLockByK(ctx context.Context, tx hcommon.DbExeAble, row *model.DBTAppLock) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_lock
SET
    v=:v,
    create_time=:create_time
WHERE
	k=:k`,
		gin.H{
			"id":          row.ID,
			"k":           row.K,
			"v":           row.V,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTProductColByName 根据id查询
func SQLGetTProductColByName(ctx context.Context, tx hcommon.DbExeAble, cols []string, appName string) (*model.DBTProduct, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product
WHERE
	app_name=:app_name`)

	var row model.DBTProduct
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"app_name": appName,
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

// SQLGetTAddressKeyColFreeForUpdate 根据id查询
func SQLGetTAddressKeyColFreeForUpdate(ctx context.Context, tx hcommon.DbExeAble, cols []string, symbol string) (*model.DBTAddressKey, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
WHERE
	use_tag=0
	AND symbol=:symbol
LIMIT 1
FOR UPDATE`)

	var row model.DBTAddressKey
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"symbol": symbol,
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

// SQLUpdateTAddressKeyUseTag 更新
func SQLUpdateTAddressKeyUseTag(ctx context.Context, tx hcommon.DbExeAble, row *model.DBTAddressKey) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_address_key
SET
    use_tag=:use_tag
WHERE
	id=:id`,
		gin.H{
			"id":      row.ID,
			"use_tag": row.UseTag,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLSelectTTxColByStatus 根据ids获取
func SQLSelectTTxColByStatus(ctx context.Context, tx hcommon.DbExeAble, cols []string, status int64) ([]*model.DBTTx, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
WHERE
	handle_status=:handle_status`)

	var rows []*model.DBTTx
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

// SQLSelectTProductNotifyColByStatusAndTime 根据ids获取
func SQLSelectTProductNotifyColByStatusAndTime(ctx context.Context, tx hcommon.DbExeAble, cols []string, status int64, t int64) ([]*model.DBTProductNotify, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
WHERE
	handle_status=:handle_status
	AND update_time<:update_time`)

	var rows []*model.DBTProductNotify
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"handle_status": status,
			"update_time":   t,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTProductNotifyStatusByID 更新
func SQLUpdateTProductNotifyStatusByID(ctx context.Context, tx hcommon.DbExeAble, row *model.DBTProductNotify) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_product_notify
SET
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    update_time=:update_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"update_time":   row.UpdateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLSelectTAppConfigTokenColAll 根据ids获取
func SQLSelectTAppConfigTokenColAll(ctx context.Context, tx hcommon.DbExeAble, cols []string) ([]*model.DBTAppConfigToken, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token`)

	var rows []*model.DBTAppConfigToken
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxErc20ColByStatus 根据ids获取
func SQLSelectTTxErc20ColByStatus(ctx context.Context, tx hcommon.DbExeAble, cols []string, status int64) ([]*model.DBTTxErc20, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
WHERE
	handle_status=:handle_status`)

	var rows []*model.DBTTxErc20
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

// SQLSelectTTxErc20ColByOrg 获取未整理交易
func SQLSelectTTxErc20ColByOrg(ctx context.Context, tx hcommon.DbExeAble, cols []string, orgStatuses []int64) ([]*model.DBTTxErc20, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
WHERE
	org_status IN (:org_status)`)

	var rows []*model.DBTTxErc20
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"org_status": orgStatuses,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxErc20OrgStatusByIDs 更新
func SQLUpdateTTxErc20OrgStatusByIDs(ctx context.Context, tx hcommon.DbExeAble, ids []int64, row model.DBTTxErc20) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_erc20
SET
    org_status=:org_status,
    org_msg=:org_msg,
    org_time=:org_time
WHERE
	id IN (:ids)`,
		gin.H{
			"ids":        ids,
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

// SQLUpdateTTxErc20OrgStatusByTxHashed 更新
func SQLUpdateTTxErc20OrgStatusByTxHashed(ctx context.Context, tx hcommon.DbExeAble, txHashes []string, row model.DBTTxErc20) (int64, error) {
	if len(txHashes) == 0 {
		return 0, nil
	}
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_erc20
SET
    org_status=:org_status,
    org_msg=:org_msg,
    org_time=:org_time
WHERE
	tx_id IN (:tx_ids)`,
		gin.H{
			"tx_ids":     txHashes,
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
