package model

import (
	"context"
	"go-dc-wallet/hcommon"
	"strings"

	"github.com/gin-gonic/gin"
)

// SQLCreateTAddressKey 创建
func SQLCreateTAddressKey(ctx context.Context, tx hcommon.DbExeAble, row *DBTAddressKey) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_address_key (
    id,
    address,
    pwd,
    use_tag
) VALUES (
    :id,
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"id":      row.ID,
				"address": row.Address,
				"pwd":     row.Pwd,
				"use_tag": row.UseTag,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_address_key (
    address,
    pwd,
    use_tag
) VALUES (
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"address": row.Address,
				"pwd":     row.Pwd,
				"use_tag": row.UseTag,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTAddressKey 创建
func SQLCreateIgnoreTAddressKey(ctx context.Context, tx hcommon.DbExeAble, row *DBTAddressKey) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_address_key (
    id,
    address,
    pwd,
    use_tag
) VALUES (
    :id,
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"id":      row.ID,
				"address": row.Address,
				"pwd":     row.Pwd,
				"use_tag": row.UseTag,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_address_key (
    address,
    pwd,
    use_tag
) VALUES (
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"address": row.Address,
				"pwd":     row.Pwd,
				"use_tag": row.UseTag,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAddressKey 创建多个
func SQLCreateManyTAddressKey(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAddressKey) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.Address,
					row.Pwd,
					row.UseTag,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.Address,
					row.Pwd,
					row.UseTag,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_address_key (
    id,
    address,
    pwd,
    use_tag
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_address_key (
    address,
    pwd,
    use_tag
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateIgnoreManyTAddressKey 创建多个
func SQLCreateIgnoreManyTAddressKey(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAddressKey) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.Address,
					row.Pwd,
					row.UseTag,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.Address,
					row.Pwd,
					row.UseTag,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_address_key (
    id,
    address,
    pwd,
    use_tag
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_address_key (
    address,
    pwd,
    use_tag
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAddressKey 根据id查询
func SQLGetTAddressKey(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAddressKey, error) {
	var row DBTAddressKey
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    address,
    pwd,
    use_tag
FROM
	t_address_key
WHERE
	id=:id`,
		gin.H{
			"id": id,
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

// SQLGetTAddressKeyCol 根据id查询
func SQLGetTAddressKeyCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAddressKey, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
WHERE
	id=:id`)

	var row DBTAddressKey
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"id": id,
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

// SQLSelectTAddressKey 根据ids获取
func SQLSelectTAddressKey(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAddressKey, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAddressKey
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    address,
    pwd,
    use_tag
FROM
	t_address_key
WHERE
	id IN (:ids)`,
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAddressKeyCol 根据ids获取
func SQLSelectTAddressKeyCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAddressKey, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
WHERE
	id IN (:ids)`)

	var rows []*DBTAddressKey
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAddressKey 更新
func SQLUpdateTAddressKey(ctx context.Context, tx hcommon.DbExeAble, row *DBTAddressKey) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_address_key
SET
    address=:address,
    pwd=:pwd,
    use_tag=:use_tag
WHERE
	id=:id`,
		gin.H{
			"id":      row.ID,
			"address": row.Address,
			"pwd":     row.Pwd,
			"use_tag": row.UseTag,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTAddressKey 删除
func SQLDeleteTAddressKey(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_address_key
WHERE
	id=:id`,
		gin.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppConfigInt 创建
func SQLCreateTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigInt) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_int (
    id,
    k,
    v
) VALUES (
    :id,
    :k,
    :v
)`,
			gin.H{
				"id": row.ID,
				"k":  row.K,
				"v":  row.V,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_int (
    k,
    v
) VALUES (
    :k,
    :v
)`,
			gin.H{
				"k": row.K,
				"v": row.V,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTAppConfigInt 创建
func SQLCreateIgnoreTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigInt) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_int (
    id,
    k,
    v
) VALUES (
    :id,
    :k,
    :v
)`,
			gin.H{
				"id": row.ID,
				"k":  row.K,
				"v":  row.V,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_int (
    k,
    v
) VALUES (
    :k,
    :v
)`,
			gin.H{
				"k": row.K,
				"v": row.V,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppConfigInt 创建多个
func SQLCreateManyTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigInt) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.K,
					row.V,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.K,
					row.V,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_int (
    id,
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_int (
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateIgnoreManyTAppConfigInt 创建多个
func SQLCreateIgnoreManyTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigInt) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.K,
					row.V,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.K,
					row.V,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_int (
    id,
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_int (
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppConfigInt 根据id查询
func SQLGetTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAppConfigInt, error) {
	var row DBTAppConfigInt
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
	id=:id`,
		gin.H{
			"id": id,
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

// SQLGetTAppConfigIntCol 根据id查询
func SQLGetTAppConfigIntCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigInt, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_int
WHERE
	id=:id`)

	var row DBTAppConfigInt
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"id": id,
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

// SQLSelectTAppConfigInt 根据ids获取
func SQLSelectTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAppConfigInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAppConfigInt
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    k,
    v
FROM
	t_app_config_int
WHERE
	id IN (:ids)`,
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppConfigIntCol 根据ids获取
func SQLSelectTAppConfigIntCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAppConfigInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_int
WHERE
	id IN (:ids)`)

	var rows []*DBTAppConfigInt
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppConfigInt 更新
func SQLUpdateTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigInt) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_int
SET
    k=:k,
    v=:v
WHERE
	id=:id`,
		gin.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTAppConfigInt 删除
func SQLDeleteTAppConfigInt(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_int
WHERE
	id=:id`,
		gin.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppStatusInt 创建
func SQLCreateTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppStatusInt) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_status_int (
    id,
    k,
    v
) VALUES (
    :id,
    :k,
    :v
)`,
			gin.H{
				"id": row.ID,
				"k":  row.K,
				"v":  row.V,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_status_int (
    k,
    v
) VALUES (
    :k,
    :v
)`,
			gin.H{
				"k": row.K,
				"v": row.V,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTAppStatusInt 创建
func SQLCreateIgnoreTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppStatusInt) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_status_int (
    id,
    k,
    v
) VALUES (
    :id,
    :k,
    :v
)`,
			gin.H{
				"id": row.ID,
				"k":  row.K,
				"v":  row.V,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_status_int (
    k,
    v
) VALUES (
    :k,
    :v
)`,
			gin.H{
				"k": row.K,
				"v": row.V,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppStatusInt 创建多个
func SQLCreateManyTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppStatusInt) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.K,
					row.V,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.K,
					row.V,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_status_int (
    id,
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_status_int (
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateIgnoreManyTAppStatusInt 创建多个
func SQLCreateIgnoreManyTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppStatusInt) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.K,
					row.V,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.K,
					row.V,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_status_int (
    id,
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_status_int (
    k,
    v
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppStatusInt 根据id查询
func SQLGetTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAppStatusInt, error) {
	var row DBTAppStatusInt
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
	id=:id`,
		gin.H{
			"id": id,
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

// SQLGetTAppStatusIntCol 根据id查询
func SQLGetTAppStatusIntCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAppStatusInt, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_status_int
WHERE
	id=:id`)

	var row DBTAppStatusInt
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"id": id,
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

// SQLSelectTAppStatusInt 根据ids获取
func SQLSelectTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAppStatusInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAppStatusInt
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    k,
    v
FROM
	t_app_status_int
WHERE
	id IN (:ids)`,
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppStatusIntCol 根据ids获取
func SQLSelectTAppStatusIntCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAppStatusInt, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_status_int
WHERE
	id IN (:ids)`)

	var rows []*DBTAppStatusInt
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppStatusInt 更新
func SQLUpdateTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppStatusInt) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_status_int
SET
    k=:k,
    v=:v
WHERE
	id=:id`,
		gin.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTAppStatusInt 删除
func SQLDeleteTAppStatusInt(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_status_int
WHERE
	id=:id`,
		gin.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTx 创建
func SQLCreateTTx(ctx context.Context, tx hcommon.DbExeAble, row *DBTTx) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx (
    id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES (
    :id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
)`,
			gin.H{
				"id":            row.ID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
				"org_status":    row.OrgStatus,
				"org_msg":       row.OrgMsg,
				"org_time":      row.OrgTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx (
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES (
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
)`,
			gin.H{
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
				"org_status":    row.OrgStatus,
				"org_msg":       row.OrgMsg,
				"org_time":      row.OrgTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTTx 创建
func SQLCreateIgnoreTTx(ctx context.Context, tx hcommon.DbExeAble, row *DBTTx) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx (
    id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES (
    :id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
)`,
			gin.H{
				"id":            row.ID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
				"org_status":    row.OrgStatus,
				"org_msg":       row.OrgMsg,
				"org_time":      row.OrgTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx (
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES (
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
)`,
			gin.H{
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
				"org_status":    row.OrgStatus,
				"org_msg":       row.OrgMsg,
				"org_time":      row.OrgTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTTx 创建多个
func SQLCreateManyTTx(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTx) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
					row.OrgStatus,
					row.OrgMsg,
					row.OrgTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
					row.OrgStatus,
					row.OrgMsg,
					row.OrgTime,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_tx (
    id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_tx (
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateIgnoreManyTTx 创建多个
func SQLCreateIgnoreManyTTx(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTx) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	var args []interface{}
	if rows[0].ID > 0 {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
					row.OrgStatus,
					row.OrgMsg,
					row.OrgTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
					row.OrgStatus,
					row.OrgMsg,
					row.OrgTime,
				},
			)
		}
	}
	var count int64
	var err error
	if rows[0].ID > 0 {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx (
    id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx (
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTx 根据id查询
func SQLGetTTx(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTTx, error) {
	var row DBTTx
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
FROM
	t_tx
WHERE
	id=:id`,
		gin.H{
			"id": id,
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

// SQLGetTTxCol 根据id查询
func SQLGetTTxCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTTx, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
WHERE
	id=:id`)

	var row DBTTx
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		gin.H{
			"id": id,
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

// SQLSelectTTx 根据ids获取
func SQLSelectTTx(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTTx, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTTx
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
FROM
	t_tx
WHERE
	id IN (:ids)`,
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxCol 根据ids获取
func SQLSelectTTxCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTTx, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
WHERE
	id IN (:ids)`)

	var rows []*DBTTx
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		gin.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTx 更新
func SQLUpdateTTx(ctx context.Context, tx hcommon.DbExeAble, row *DBTTx) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx
SET
    tx_id=:tx_id,
    from_address=:from_address,
    to_address=:to_address,
    balance=:balance,
    balance_real=:balance_real,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time,
    org_status=:org_status,
    org_msg=:org_msg,
    org_time=:org_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"balance":       row.Balance,
			"balance_real":  row.BalanceReal,
			"create_time":   row.CreateTime,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_time":   row.HandleTime,
			"org_status":    row.OrgStatus,
			"org_msg":       row.OrgMsg,
			"org_time":      row.OrgTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTTx 删除
func SQLDeleteTTx(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx
WHERE
	id=:id`,
		gin.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}
