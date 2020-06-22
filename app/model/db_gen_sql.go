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
    symbol,
    address,
    pwd,
    use_tag
) VALUES (
    :id,
    :symbol,
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"id":      row.ID,
				"symbol":  row.Symbol,
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
    symbol,
    address,
    pwd,
    use_tag
) VALUES (
    :symbol,
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"symbol":  row.Symbol,
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
    symbol,
    address,
    pwd,
    use_tag
) VALUES (
    :id,
    :symbol,
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"id":      row.ID,
				"symbol":  row.Symbol,
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
    symbol,
    address,
    pwd,
    use_tag
) VALUES (
    :symbol,
    :address,
    :pwd,
    :use_tag
)`,
			gin.H{
				"symbol":  row.Symbol,
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
					row.Symbol,
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
					row.Symbol,
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
    symbol,
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
    symbol,
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
					row.Symbol,
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
					row.Symbol,
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
    symbol,
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
    symbol,
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
    symbol,
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
    symbol,
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
    symbol=:symbol,
    address=:address,
    pwd=:pwd,
    use_tag=:use_tag
WHERE
	id=:id`,
		gin.H{
			"id":      row.ID,
			"symbol":  row.Symbol,
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

// SQLCreateTAppConfigStr 创建
func SQLCreateTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigStr) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_str (
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
			`INSERT INTO t_app_config_str (
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

// SQLCreateIgnoreTAppConfigStr 创建
func SQLCreateIgnoreTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigStr) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_str (
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
			`INSERT IGNORE INTO t_app_config_str (
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

// SQLCreateManyTAppConfigStr 创建多个
func SQLCreateManyTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigStr) (int64, error) {
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
			`INSERT INTO t_app_config_str (
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
			`INSERT INTO t_app_config_str (
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

// SQLCreateIgnoreManyTAppConfigStr 创建多个
func SQLCreateIgnoreManyTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigStr) (int64, error) {
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
			`INSERT IGNORE INTO t_app_config_str (
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
			`INSERT IGNORE INTO t_app_config_str (
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

// SQLGetTAppConfigStr 根据id查询
func SQLGetTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAppConfigStr, error) {
	var row DBTAppConfigStr
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

// SQLGetTAppConfigStrCol 根据id查询
func SQLGetTAppConfigStrCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigStr, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_str
WHERE
	id=:id`)

	var row DBTAppConfigStr
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

// SQLSelectTAppConfigStr 根据ids获取
func SQLSelectTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAppConfigStr, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAppConfigStr
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    k,
    v
FROM
	t_app_config_str
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

// SQLSelectTAppConfigStrCol 根据ids获取
func SQLSelectTAppConfigStrCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAppConfigStr, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_str
WHERE
	id IN (:ids)`)

	var rows []*DBTAppConfigStr
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

// SQLUpdateTAppConfigStr 更新
func SQLUpdateTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigStr) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_str
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

// SQLDeleteTAppConfigStr 删除
func SQLDeleteTAppConfigStr(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_str
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

// SQLCreateTAppConfigToken 创建
func SQLCreateTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigToken) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_token (
    id,
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES (
    :id,
    :token_address,
    :token_decimals,
    :token_symbol,
    :cold_address,
    :hot_address,
    :org_min_balance,
    :create_time
)`,
			gin.H{
				"id":              row.ID,
				"token_address":   row.TokenAddress,
				"token_decimals":  row.TokenDecimals,
				"token_symbol":    row.TokenSymbol,
				"cold_address":    row.ColdAddress,
				"hot_address":     row.HotAddress,
				"org_min_balance": row.OrgMinBalance,
				"create_time":     row.CreateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_token (
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES (
    :token_address,
    :token_decimals,
    :token_symbol,
    :cold_address,
    :hot_address,
    :org_min_balance,
    :create_time
)`,
			gin.H{
				"token_address":   row.TokenAddress,
				"token_decimals":  row.TokenDecimals,
				"token_symbol":    row.TokenSymbol,
				"cold_address":    row.ColdAddress,
				"hot_address":     row.HotAddress,
				"org_min_balance": row.OrgMinBalance,
				"create_time":     row.CreateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTAppConfigToken 创建
func SQLCreateIgnoreTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigToken) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_token (
    id,
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES (
    :id,
    :token_address,
    :token_decimals,
    :token_symbol,
    :cold_address,
    :hot_address,
    :org_min_balance,
    :create_time
)`,
			gin.H{
				"id":              row.ID,
				"token_address":   row.TokenAddress,
				"token_decimals":  row.TokenDecimals,
				"token_symbol":    row.TokenSymbol,
				"cold_address":    row.ColdAddress,
				"hot_address":     row.HotAddress,
				"org_min_balance": row.OrgMinBalance,
				"create_time":     row.CreateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_token (
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES (
    :token_address,
    :token_decimals,
    :token_symbol,
    :cold_address,
    :hot_address,
    :org_min_balance,
    :create_time
)`,
			gin.H{
				"token_address":   row.TokenAddress,
				"token_decimals":  row.TokenDecimals,
				"token_symbol":    row.TokenSymbol,
				"cold_address":    row.ColdAddress,
				"hot_address":     row.HotAddress,
				"org_min_balance": row.OrgMinBalance,
				"create_time":     row.CreateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppConfigToken 创建多个
func SQLCreateManyTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigToken) (int64, error) {
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
					row.TokenAddress,
					row.TokenDecimals,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.OrgMinBalance,
					row.CreateTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TokenAddress,
					row.TokenDecimals,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.OrgMinBalance,
					row.CreateTime,
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
			`INSERT INTO t_app_config_token (
    id,
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_token (
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
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

// SQLCreateIgnoreManyTAppConfigToken 创建多个
func SQLCreateIgnoreManyTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigToken) (int64, error) {
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
					row.TokenAddress,
					row.TokenDecimals,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.OrgMinBalance,
					row.CreateTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TokenAddress,
					row.TokenDecimals,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.OrgMinBalance,
					row.CreateTime,
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
			`INSERT IGNORE INTO t_app_config_token (
    id,
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_token (
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
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

// SQLGetTAppConfigToken 根据id查询
func SQLGetTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAppConfigToken, error) {
	var row DBTAppConfigToken
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
FROM
	t_app_config_token
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

// SQLGetTAppConfigTokenCol 根据id查询
func SQLGetTAppConfigTokenCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigToken, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token
WHERE
	id=:id`)

	var row DBTAppConfigToken
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

// SQLSelectTAppConfigToken 根据ids获取
func SQLSelectTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAppConfigToken, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAppConfigToken
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
FROM
	t_app_config_token
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

// SQLSelectTAppConfigTokenCol 根据ids获取
func SQLSelectTAppConfigTokenCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAppConfigToken, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token
WHERE
	id IN (:ids)`)

	var rows []*DBTAppConfigToken
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

// SQLUpdateTAppConfigToken 更新
func SQLUpdateTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigToken) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_token
SET
    token_address=:token_address,
    token_decimals=:token_decimals,
    token_symbol=:token_symbol,
    cold_address=:cold_address,
    hot_address=:hot_address,
    org_min_balance=:org_min_balance,
    create_time=:create_time
WHERE
	id=:id`,
		gin.H{
			"id":              row.ID,
			"token_address":   row.TokenAddress,
			"token_decimals":  row.TokenDecimals,
			"token_symbol":    row.TokenSymbol,
			"cold_address":    row.ColdAddress,
			"hot_address":     row.HotAddress,
			"org_min_balance": row.OrgMinBalance,
			"create_time":     row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTAppConfigToken 删除
func SQLDeleteTAppConfigToken(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_token
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

// SQLCreateTAppConfigTokenBtc 创建
func SQLCreateTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigTokenBtc) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_token_btc (
    id,
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
) VALUES (
    :id,
    :token_index,
    :token_symbol,
    :cold_address,
    :hot_address,
    :tx_org_min_balance,
    :create_at
)`,
			gin.H{
				"id":                 row.ID,
				"token_index":        row.TokenIndex,
				"token_symbol":       row.TokenSymbol,
				"cold_address":       row.ColdAddress,
				"hot_address":        row.HotAddress,
				"tx_org_min_balance": row.TxOrgMinBalance,
				"create_at":          row.CreateAt,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_token_btc (
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
) VALUES (
    :token_index,
    :token_symbol,
    :cold_address,
    :hot_address,
    :tx_org_min_balance,
    :create_at
)`,
			gin.H{
				"token_index":        row.TokenIndex,
				"token_symbol":       row.TokenSymbol,
				"cold_address":       row.ColdAddress,
				"hot_address":        row.HotAddress,
				"tx_org_min_balance": row.TxOrgMinBalance,
				"create_at":          row.CreateAt,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTAppConfigTokenBtc 创建
func SQLCreateIgnoreTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigTokenBtc) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_token_btc (
    id,
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
) VALUES (
    :id,
    :token_index,
    :token_symbol,
    :cold_address,
    :hot_address,
    :tx_org_min_balance,
    :create_at
)`,
			gin.H{
				"id":                 row.ID,
				"token_index":        row.TokenIndex,
				"token_symbol":       row.TokenSymbol,
				"cold_address":       row.ColdAddress,
				"hot_address":        row.HotAddress,
				"tx_org_min_balance": row.TxOrgMinBalance,
				"create_at":          row.CreateAt,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_token_btc (
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
) VALUES (
    :token_index,
    :token_symbol,
    :cold_address,
    :hot_address,
    :tx_org_min_balance,
    :create_at
)`,
			gin.H{
				"token_index":        row.TokenIndex,
				"token_symbol":       row.TokenSymbol,
				"cold_address":       row.ColdAddress,
				"hot_address":        row.HotAddress,
				"tx_org_min_balance": row.TxOrgMinBalance,
				"create_at":          row.CreateAt,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppConfigTokenBtc 创建多个
func SQLCreateManyTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigTokenBtc) (int64, error) {
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
					row.TokenIndex,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.TxOrgMinBalance,
					row.CreateAt,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TokenIndex,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.TxOrgMinBalance,
					row.CreateAt,
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
			`INSERT INTO t_app_config_token_btc (
    id,
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_config_token_btc (
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
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

// SQLCreateIgnoreManyTAppConfigTokenBtc 创建多个
func SQLCreateIgnoreManyTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppConfigTokenBtc) (int64, error) {
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
					row.TokenIndex,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.TxOrgMinBalance,
					row.CreateAt,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TokenIndex,
					row.TokenSymbol,
					row.ColdAddress,
					row.HotAddress,
					row.TxOrgMinBalance,
					row.CreateAt,
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
			`INSERT IGNORE INTO t_app_config_token_btc (
    id,
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_config_token_btc (
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
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

// SQLGetTAppConfigTokenBtc 根据id查询
func SQLGetTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAppConfigTokenBtc, error) {
	var row DBTAppConfigTokenBtc
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
FROM
	t_app_config_token_btc
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

// SQLGetTAppConfigTokenBtcCol 根据id查询
func SQLGetTAppConfigTokenBtcCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigTokenBtc, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token_btc
WHERE
	id=:id`)

	var row DBTAppConfigTokenBtc
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

// SQLSelectTAppConfigTokenBtc 根据ids获取
func SQLSelectTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAppConfigTokenBtc, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAppConfigTokenBtc
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    tx_org_min_balance,
    create_at
FROM
	t_app_config_token_btc
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

// SQLSelectTAppConfigTokenBtcCol 根据ids获取
func SQLSelectTAppConfigTokenBtcCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAppConfigTokenBtc, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token_btc
WHERE
	id IN (:ids)`)

	var rows []*DBTAppConfigTokenBtc
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

// SQLUpdateTAppConfigTokenBtc 更新
func SQLUpdateTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppConfigTokenBtc) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_token_btc
SET
    token_index=:token_index,
    token_symbol=:token_symbol,
    cold_address=:cold_address,
    hot_address=:hot_address,
    tx_org_min_balance=:tx_org_min_balance,
    create_at=:create_at
WHERE
	id=:id`,
		gin.H{
			"id":                 row.ID,
			"token_index":        row.TokenIndex,
			"token_symbol":       row.TokenSymbol,
			"cold_address":       row.ColdAddress,
			"hot_address":        row.HotAddress,
			"tx_org_min_balance": row.TxOrgMinBalance,
			"create_at":          row.CreateAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTAppConfigTokenBtc 删除
func SQLDeleteTAppConfigTokenBtc(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_token_btc
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

// SQLCreateTAppLock 创建
func SQLCreateTAppLock(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppLock) (int64, error) {
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
)`,
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
)`,
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

// SQLCreateIgnoreTAppLock 创建
func SQLCreateIgnoreTAppLock(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppLock) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_lock (
    id,
    k,
    v,
    create_time
) VALUES (
    :id,
    :k,
    :v,
    :create_time
)`,
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
			`INSERT IGNORE INTO t_app_lock (
    k,
    v,
    create_time
) VALUES (
    :k,
    :v,
    :create_time
)`,
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

// SQLCreateManyTAppLock 创建多个
func SQLCreateManyTAppLock(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppLock) (int64, error) {
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
					row.CreateTime,
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
					row.CreateTime,
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
			`INSERT INTO t_app_lock (
    id,
    k,
    v,
    create_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_app_lock (
    k,
    v,
    create_time
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

// SQLCreateIgnoreManyTAppLock 创建多个
func SQLCreateIgnoreManyTAppLock(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTAppLock) (int64, error) {
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
					row.CreateTime,
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
					row.CreateTime,
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
			`INSERT IGNORE INTO t_app_lock (
    id,
    k,
    v,
    create_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_app_lock (
    k,
    v,
    create_time
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

// SQLGetTAppLock 根据id查询
func SQLGetTAppLock(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTAppLock, error) {
	var row DBTAppLock
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    k,
    v,
    create_time
FROM
	t_app_lock
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

// SQLGetTAppLockCol 根据id查询
func SQLGetTAppLockCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTAppLock, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
WHERE
	id=:id`)

	var row DBTAppLock
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

// SQLSelectTAppLock 根据ids获取
func SQLSelectTAppLock(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTAppLock, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTAppLock
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    k,
    v,
    create_time
FROM
	t_app_lock
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

// SQLSelectTAppLockCol 根据ids获取
func SQLSelectTAppLockCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTAppLock, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
WHERE
	id IN (:ids)`)

	var rows []*DBTAppLock
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

// SQLUpdateTAppLock 更新
func SQLUpdateTAppLock(ctx context.Context, tx hcommon.DbExeAble, row *DBTAppLock) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_lock
SET
    k=:k,
    v=:v,
    create_time=:create_time
WHERE
	id=:id`,
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

// SQLDeleteTAppLock 删除
func SQLDeleteTAppLock(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_lock
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

// SQLCreateTProduct 创建
func SQLCreateTProduct(ctx context.Context, tx hcommon.DbExeAble, row *DBTProduct) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_product (
    id,
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES (
    :id,
    :app_name,
    :app_sk,
    :cb_url,
    :whitelist_ip
)`,
			gin.H{
				"id":           row.ID,
				"app_name":     row.AppName,
				"app_sk":       row.AppSk,
				"cb_url":       row.CbURL,
				"whitelist_ip": row.WhitelistIP,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_product (
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES (
    :app_name,
    :app_sk,
    :cb_url,
    :whitelist_ip
)`,
			gin.H{
				"app_name":     row.AppName,
				"app_sk":       row.AppSk,
				"cb_url":       row.CbURL,
				"whitelist_ip": row.WhitelistIP,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTProduct 创建
func SQLCreateIgnoreTProduct(ctx context.Context, tx hcommon.DbExeAble, row *DBTProduct) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product (
    id,
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES (
    :id,
    :app_name,
    :app_sk,
    :cb_url,
    :whitelist_ip
)`,
			gin.H{
				"id":           row.ID,
				"app_name":     row.AppName,
				"app_sk":       row.AppSk,
				"cb_url":       row.CbURL,
				"whitelist_ip": row.WhitelistIP,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product (
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES (
    :app_name,
    :app_sk,
    :cb_url,
    :whitelist_ip
)`,
			gin.H{
				"app_name":     row.AppName,
				"app_sk":       row.AppSk,
				"cb_url":       row.CbURL,
				"whitelist_ip": row.WhitelistIP,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTProduct 创建多个
func SQLCreateManyTProduct(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTProduct) (int64, error) {
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
					row.AppName,
					row.AppSk,
					row.CbURL,
					row.WhitelistIP,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.AppName,
					row.AppSk,
					row.CbURL,
					row.WhitelistIP,
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
			`INSERT INTO t_product (
    id,
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_product (
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
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

// SQLCreateIgnoreManyTProduct 创建多个
func SQLCreateIgnoreManyTProduct(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTProduct) (int64, error) {
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
					row.AppName,
					row.AppSk,
					row.CbURL,
					row.WhitelistIP,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.AppName,
					row.AppSk,
					row.CbURL,
					row.WhitelistIP,
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
			`INSERT IGNORE INTO t_product (
    id,
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product (
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
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

// SQLGetTProduct 根据id查询
func SQLGetTProduct(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTProduct, error) {
	var row DBTProduct
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
FROM
	t_product
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

// SQLGetTProductCol 根据id查询
func SQLGetTProductCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTProduct, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product
WHERE
	id=:id`)

	var row DBTProduct
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

// SQLSelectTProduct 根据ids获取
func SQLSelectTProduct(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTProduct, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTProduct
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
FROM
	t_product
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

// SQLSelectTProductCol 根据ids获取
func SQLSelectTProductCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTProduct, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product
WHERE
	id IN (:ids)`)

	var rows []*DBTProduct
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

// SQLUpdateTProduct 更新
func SQLUpdateTProduct(ctx context.Context, tx hcommon.DbExeAble, row *DBTProduct) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_product
SET
    app_name=:app_name,
    app_sk=:app_sk,
    cb_url=:cb_url,
    whitelist_ip=:whitelist_ip
WHERE
	id=:id`,
		gin.H{
			"id":           row.ID,
			"app_name":     row.AppName,
			"app_sk":       row.AppSk,
			"cb_url":       row.CbURL,
			"whitelist_ip": row.WhitelistIP,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTProduct 删除
func SQLDeleteTProduct(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_product
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

// SQLCreateTProductNonce 创建
func SQLCreateTProductNonce(ctx context.Context, tx hcommon.DbExeAble, row *DBTProductNonce) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_product_nonce (
    id,
    c,
    create_time
) VALUES (
    :id,
    :c,
    :create_time
)`,
			gin.H{
				"id":          row.ID,
				"c":           row.C,
				"create_time": row.CreateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_product_nonce (
    c,
    create_time
) VALUES (
    :c,
    :create_time
)`,
			gin.H{
				"c":           row.C,
				"create_time": row.CreateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTProductNonce 创建
func SQLCreateIgnoreTProductNonce(ctx context.Context, tx hcommon.DbExeAble, row *DBTProductNonce) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product_nonce (
    id,
    c,
    create_time
) VALUES (
    :id,
    :c,
    :create_time
)`,
			gin.H{
				"id":          row.ID,
				"c":           row.C,
				"create_time": row.CreateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product_nonce (
    c,
    create_time
) VALUES (
    :c,
    :create_time
)`,
			gin.H{
				"c":           row.C,
				"create_time": row.CreateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTProductNonce 创建多个
func SQLCreateManyTProductNonce(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTProductNonce) (int64, error) {
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
					row.C,
					row.CreateTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.C,
					row.CreateTime,
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
			`INSERT INTO t_product_nonce (
    id,
    c,
    create_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_product_nonce (
    c,
    create_time
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

// SQLCreateIgnoreManyTProductNonce 创建多个
func SQLCreateIgnoreManyTProductNonce(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTProductNonce) (int64, error) {
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
					row.C,
					row.CreateTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.C,
					row.CreateTime,
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
			`INSERT IGNORE INTO t_product_nonce (
    id,
    c,
    create_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product_nonce (
    c,
    create_time
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

// SQLGetTProductNonce 根据id查询
func SQLGetTProductNonce(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTProductNonce, error) {
	var row DBTProductNonce
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    c,
    create_time
FROM
	t_product_nonce
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

// SQLGetTProductNonceCol 根据id查询
func SQLGetTProductNonceCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTProductNonce, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_nonce
WHERE
	id=:id`)

	var row DBTProductNonce
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

// SQLSelectTProductNonce 根据ids获取
func SQLSelectTProductNonce(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTProductNonce, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTProductNonce
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    c,
    create_time
FROM
	t_product_nonce
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

// SQLSelectTProductNonceCol 根据ids获取
func SQLSelectTProductNonceCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTProductNonce, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_nonce
WHERE
	id IN (:ids)`)

	var rows []*DBTProductNonce
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

// SQLUpdateTProductNonce 更新
func SQLUpdateTProductNonce(ctx context.Context, tx hcommon.DbExeAble, row *DBTProductNonce) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_product_nonce
SET
    c=:c,
    create_time=:create_time
WHERE
	id=:id`,
		gin.H{
			"id":          row.ID,
			"c":           row.C,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTProductNonce 删除
func SQLDeleteTProductNonce(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_product_nonce
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

// SQLCreateTProductNotify 创建
func SQLCreateTProductNotify(ctx context.Context, tx hcommon.DbExeAble, row *DBTProductNotify) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_product_notify (
    id,
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES (
    :id,
    :nonce,
    :product_id,
    :item_type,
    :item_id,
    :notify_type,
    :url,
    :msg,
    :handle_status,
    :handle_msg,
    :create_time,
    :update_time
)`,
			gin.H{
				"id":            row.ID,
				"nonce":         row.Nonce,
				"product_id":    row.ProductID,
				"item_type":     row.ItemType,
				"item_id":       row.ItemID,
				"notify_type":   row.NotifyType,
				"url":           row.URL,
				"msg":           row.Msg,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"create_time":   row.CreateTime,
				"update_time":   row.UpdateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_product_notify (
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES (
    :nonce,
    :product_id,
    :item_type,
    :item_id,
    :notify_type,
    :url,
    :msg,
    :handle_status,
    :handle_msg,
    :create_time,
    :update_time
)`,
			gin.H{
				"nonce":         row.Nonce,
				"product_id":    row.ProductID,
				"item_type":     row.ItemType,
				"item_id":       row.ItemID,
				"notify_type":   row.NotifyType,
				"url":           row.URL,
				"msg":           row.Msg,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"create_time":   row.CreateTime,
				"update_time":   row.UpdateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTProductNotify 创建
func SQLCreateIgnoreTProductNotify(ctx context.Context, tx hcommon.DbExeAble, row *DBTProductNotify) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product_notify (
    id,
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES (
    :id,
    :nonce,
    :product_id,
    :item_type,
    :item_id,
    :notify_type,
    :url,
    :msg,
    :handle_status,
    :handle_msg,
    :create_time,
    :update_time
)`,
			gin.H{
				"id":            row.ID,
				"nonce":         row.Nonce,
				"product_id":    row.ProductID,
				"item_type":     row.ItemType,
				"item_id":       row.ItemID,
				"notify_type":   row.NotifyType,
				"url":           row.URL,
				"msg":           row.Msg,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"create_time":   row.CreateTime,
				"update_time":   row.UpdateTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product_notify (
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES (
    :nonce,
    :product_id,
    :item_type,
    :item_id,
    :notify_type,
    :url,
    :msg,
    :handle_status,
    :handle_msg,
    :create_time,
    :update_time
)`,
			gin.H{
				"nonce":         row.Nonce,
				"product_id":    row.ProductID,
				"item_type":     row.ItemType,
				"item_id":       row.ItemID,
				"notify_type":   row.NotifyType,
				"url":           row.URL,
				"msg":           row.Msg,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"create_time":   row.CreateTime,
				"update_time":   row.UpdateTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTProductNotify 创建多个
func SQLCreateManyTProductNotify(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTProductNotify) (int64, error) {
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
					row.Nonce,
					row.ProductID,
					row.ItemType,
					row.ItemID,
					row.NotifyType,
					row.URL,
					row.Msg,
					row.HandleStatus,
					row.HandleMsg,
					row.CreateTime,
					row.UpdateTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.Nonce,
					row.ProductID,
					row.ItemType,
					row.ItemID,
					row.NotifyType,
					row.URL,
					row.Msg,
					row.HandleStatus,
					row.HandleMsg,
					row.CreateTime,
					row.UpdateTime,
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
			`INSERT INTO t_product_notify (
    id,
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_product_notify (
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
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

// SQLCreateIgnoreManyTProductNotify 创建多个
func SQLCreateIgnoreManyTProductNotify(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTProductNotify) (int64, error) {
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
					row.Nonce,
					row.ProductID,
					row.ItemType,
					row.ItemID,
					row.NotifyType,
					row.URL,
					row.Msg,
					row.HandleStatus,
					row.HandleMsg,
					row.CreateTime,
					row.UpdateTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.Nonce,
					row.ProductID,
					row.ItemType,
					row.ItemID,
					row.NotifyType,
					row.URL,
					row.Msg,
					row.HandleStatus,
					row.HandleMsg,
					row.CreateTime,
					row.UpdateTime,
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
			`INSERT IGNORE INTO t_product_notify (
    id,
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_product_notify (
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
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

// SQLGetTProductNotify 根据id查询
func SQLGetTProductNotify(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTProductNotify, error) {
	var row DBTProductNotify
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
FROM
	t_product_notify
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

// SQLGetTProductNotifyCol 根据id查询
func SQLGetTProductNotifyCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTProductNotify, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
WHERE
	id=:id`)

	var row DBTProductNotify
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

// SQLSelectTProductNotify 根据ids获取
func SQLSelectTProductNotify(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTProductNotify, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTProductNotify
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
FROM
	t_product_notify
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

// SQLSelectTProductNotifyCol 根据ids获取
func SQLSelectTProductNotifyCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTProductNotify, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
WHERE
	id IN (:ids)`)

	var rows []*DBTProductNotify
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

// SQLUpdateTProductNotify 更新
func SQLUpdateTProductNotify(ctx context.Context, tx hcommon.DbExeAble, row *DBTProductNotify) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_product_notify
SET
    nonce=:nonce,
    product_id=:product_id,
    item_type=:item_type,
    item_id=:item_id,
    notify_type=:notify_type,
    url=:url,
    msg=:msg,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    create_time=:create_time,
    update_time=:update_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"nonce":         row.Nonce,
			"product_id":    row.ProductID,
			"item_type":     row.ItemType,
			"item_id":       row.ItemID,
			"notify_type":   row.NotifyType,
			"url":           row.URL,
			"msg":           row.Msg,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"create_time":   row.CreateTime,
			"update_time":   row.UpdateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTProductNotify 删除
func SQLDeleteTProductNotify(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_product_notify
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

// SQLCreateTSend 创建
func SQLCreateTSend(ctx context.Context, tx hcommon.DbExeAble, row *DBTSend) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_send (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :nonce,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"nonce":         row.Nonce,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_send (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :nonce,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"nonce":         row.Nonce,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTSend 创建
func SQLCreateIgnoreTSend(ctx context.Context, tx hcommon.DbExeAble, row *DBTSend) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_send (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :nonce,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"nonce":         row.Nonce,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_send (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :nonce,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"nonce":         row.Nonce,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTSend 创建多个
func SQLCreateManyTSend(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTSend) (int64, error) {
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
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Nonce,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Nonce,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT INTO t_send (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_send (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLCreateIgnoreManyTSend 创建多个
func SQLCreateIgnoreManyTSend(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTSend) (int64, error) {
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
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Nonce,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Nonce,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT IGNORE INTO t_send (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_send (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLGetTSend 根据id查询
func SQLGetTSend(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTSend, error) {
	var row DBTSend
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_send
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

// SQLGetTSendCol 根据id查询
func SQLGetTSendCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTSend, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send
WHERE
	id=:id`)

	var row DBTSend
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

// SQLSelectTSend 根据ids获取
func SQLSelectTSend(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTSend, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTSend
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    nonce,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_send
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

// SQLSelectTSendCol 根据ids获取
func SQLSelectTSendCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTSend, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send
WHERE
	id IN (:ids)`)

	var rows []*DBTSend
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

// SQLUpdateTSend 更新
func SQLUpdateTSend(ctx context.Context, tx hcommon.DbExeAble, row *DBTSend) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_send
SET
    related_type=:related_type,
    related_id=:related_id,
    token_id=:token_id,
    tx_id=:tx_id,
    from_address=:from_address,
    to_address=:to_address,
    balance=:balance,
    balance_real=:balance_real,
    gas=:gas,
    gas_price=:gas_price,
    nonce=:nonce,
    hex=:hex,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"balance":       row.Balance,
			"balance_real":  row.BalanceReal,
			"gas":           row.Gas,
			"gas_price":     row.GasPrice,
			"nonce":         row.Nonce,
			"hex":           row.Hex,
			"create_time":   row.CreateTime,
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

// SQLDeleteTSend 删除
func SQLDeleteTSend(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_send
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

// SQLCreateTSendBtc 创建
func SQLCreateTSendBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTSendBtc) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_send_btc (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_send_btc (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTSendBtc 创建
func SQLCreateIgnoreTSendBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTSendBtc) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_send_btc (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_send_btc (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance,
    :balance_real,
    :gas,
    :gas_price,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"related_type":  row.RelatedType,
				"related_id":    row.RelatedID,
				"token_id":      row.TokenID,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"balance":       row.Balance,
				"balance_real":  row.BalanceReal,
				"gas":           row.Gas,
				"gas_price":     row.GasPrice,
				"hex":           row.Hex,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTSendBtc 创建多个
func SQLCreateManyTSendBtc(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTSendBtc) (int64, error) {
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
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT INTO t_send_btc (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_send_btc (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLCreateIgnoreManyTSendBtc 创建多个
func SQLCreateIgnoreManyTSendBtc(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTSendBtc) (int64, error) {
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
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.RelatedType,
					row.RelatedID,
					row.TokenID,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Balance,
					row.BalanceReal,
					row.Gas,
					row.GasPrice,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT IGNORE INTO t_send_btc (
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_send_btc (
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLGetTSendBtc 根据id查询
func SQLGetTSendBtc(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTSendBtc, error) {
	var row DBTSendBtc
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_send_btc
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

// SQLGetTSendBtcCol 根据id查询
func SQLGetTSendBtcCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTSendBtc, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_btc
WHERE
	id=:id`)

	var row DBTSendBtc
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

// SQLSelectTSendBtc 根据ids获取
func SQLSelectTSendBtc(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTSendBtc, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTSendBtc
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_send_btc
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

// SQLSelectTSendBtcCol 根据ids获取
func SQLSelectTSendBtcCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTSendBtc, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_btc
WHERE
	id IN (:ids)`)

	var rows []*DBTSendBtc
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

// SQLUpdateTSendBtc 更新
func SQLUpdateTSendBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTSendBtc) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_send_btc
SET
    related_type=:related_type,
    related_id=:related_id,
    token_id=:token_id,
    tx_id=:tx_id,
    from_address=:from_address,
    to_address=:to_address,
    balance=:balance,
    balance_real=:balance_real,
    gas=:gas,
    gas_price=:gas_price,
    hex=:hex,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"balance":       row.Balance,
			"balance_real":  row.BalanceReal,
			"gas":           row.Gas,
			"gas_price":     row.GasPrice,
			"hex":           row.Hex,
			"create_time":   row.CreateTime,
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

// SQLDeleteTSendBtc 删除
func SQLDeleteTSendBtc(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_send_btc
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
    product_id,
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
    :product_id,
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
				"product_id":    row.ProductID,
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
    product_id,
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
    :product_id,
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
				"product_id":    row.ProductID,
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
    product_id,
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
    :product_id,
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
				"product_id":    row.ProductID,
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
    product_id,
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
    :product_id,
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
				"product_id":    row.ProductID,
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
					row.ProductID,
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
					row.ProductID,
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
    product_id,
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
    product_id,
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
					row.ProductID,
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
					row.ProductID,
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
    product_id,
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
    product_id,
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
    product_id,
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
    product_id,
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
    product_id=:product_id,
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
			"product_id":    row.ProductID,
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

// SQLCreateTTxBtc 创建
func SQLCreateTTxBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtc) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc (
    id,
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :product_id,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"product_id":    row.ProductID,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc (
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :product_id,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"product_id":    row.ProductID,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTTxBtc 创建
func SQLCreateIgnoreTTxBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtc) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc (
    id,
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :product_id,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"product_id":    row.ProductID,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc (
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :product_id,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"product_id":    row.ProductID,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTTxBtc 创建多个
func SQLCreateManyTTxBtc(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxBtc) (int64, error) {
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
					row.ProductID,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ProductID,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT INTO t_tx_btc (
    id,
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc (
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLCreateIgnoreManyTTxBtc 创建多个
func SQLCreateIgnoreManyTTxBtc(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxBtc) (int64, error) {
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
					row.ProductID,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ProductID,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT IGNORE INTO t_tx_btc (
    id,
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc (
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLGetTTxBtc 根据id查询
func SQLGetTTxBtc(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTTxBtc, error) {
	var row DBTTxBtc
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_tx_btc
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

// SQLGetTTxBtcCol 根据id查询
func SQLGetTTxBtcCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTTxBtc, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc
WHERE
	id=:id`)

	var row DBTTxBtc
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

// SQLSelectTTxBtc 根据ids获取
func SQLSelectTTxBtc(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTTxBtc, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTTxBtc
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    product_id,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_tx_btc
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

// SQLSelectTTxBtcCol 根据ids获取
func SQLSelectTTxBtcCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTTxBtc, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc
WHERE
	id IN (:ids)`)

	var rows []*DBTTxBtc
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

// SQLUpdateTTxBtc 更新
func SQLUpdateTTxBtc(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtc) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_btc
SET
    product_id=:product_id,
    block_hash=:block_hash,
    tx_id=:tx_id,
    vout_n=:vout_n,
    vout_address=:vout_address,
    vout_value=:vout_value,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"block_hash":    row.BlockHash,
			"tx_id":         row.TxID,
			"vout_n":        row.VoutN,
			"vout_address":  row.VoutAddress,
			"vout_value":    row.VoutValue,
			"create_time":   row.CreateTime,
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

// SQLDeleteTTxBtc 删除
func SQLDeleteTTxBtc(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_btc
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

// SQLCreateTTxBtcToken 创建
func SQLCreateTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtcToken) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc_token (
    id,
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES (
    :id,
    :token_index,
    :token_symbol,
    :block_hash,
    :tx_id,
    :from_address,
    :to_address,
    :value,
    :blocktime,
    :create_at,
    :handle_status,
    :handle_msg,
    :handle_at
)`,
			gin.H{
				"id":            row.ID,
				"token_index":   row.TokenIndex,
				"token_symbol":  row.TokenSymbol,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"value":         row.Value,
				"blocktime":     row.Blocktime,
				"create_at":     row.CreateAt,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_at":     row.HandleAt,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc_token (
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES (
    :token_index,
    :token_symbol,
    :block_hash,
    :tx_id,
    :from_address,
    :to_address,
    :value,
    :blocktime,
    :create_at,
    :handle_status,
    :handle_msg,
    :handle_at
)`,
			gin.H{
				"token_index":   row.TokenIndex,
				"token_symbol":  row.TokenSymbol,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"value":         row.Value,
				"blocktime":     row.Blocktime,
				"create_at":     row.CreateAt,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_at":     row.HandleAt,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTTxBtcToken 创建
func SQLCreateIgnoreTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtcToken) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc_token (
    id,
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES (
    :id,
    :token_index,
    :token_symbol,
    :block_hash,
    :tx_id,
    :from_address,
    :to_address,
    :value,
    :blocktime,
    :create_at,
    :handle_status,
    :handle_msg,
    :handle_at
)`,
			gin.H{
				"id":            row.ID,
				"token_index":   row.TokenIndex,
				"token_symbol":  row.TokenSymbol,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"value":         row.Value,
				"blocktime":     row.Blocktime,
				"create_at":     row.CreateAt,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_at":     row.HandleAt,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc_token (
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES (
    :token_index,
    :token_symbol,
    :block_hash,
    :tx_id,
    :from_address,
    :to_address,
    :value,
    :blocktime,
    :create_at,
    :handle_status,
    :handle_msg,
    :handle_at
)`,
			gin.H{
				"token_index":   row.TokenIndex,
				"token_symbol":  row.TokenSymbol,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"from_address":  row.FromAddress,
				"to_address":    row.ToAddress,
				"value":         row.Value,
				"blocktime":     row.Blocktime,
				"create_at":     row.CreateAt,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_at":     row.HandleAt,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTTxBtcToken 创建多个
func SQLCreateManyTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxBtcToken) (int64, error) {
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
					row.TokenIndex,
					row.TokenSymbol,
					row.BlockHash,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Value,
					row.Blocktime,
					row.CreateAt,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleAt,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TokenIndex,
					row.TokenSymbol,
					row.BlockHash,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Value,
					row.Blocktime,
					row.CreateAt,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleAt,
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
			`INSERT INTO t_tx_btc_token (
    id,
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc_token (
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
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

// SQLCreateIgnoreManyTTxBtcToken 创建多个
func SQLCreateIgnoreManyTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxBtcToken) (int64, error) {
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
					row.TokenIndex,
					row.TokenSymbol,
					row.BlockHash,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Value,
					row.Blocktime,
					row.CreateAt,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleAt,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.TokenIndex,
					row.TokenSymbol,
					row.BlockHash,
					row.TxID,
					row.FromAddress,
					row.ToAddress,
					row.Value,
					row.Blocktime,
					row.CreateAt,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleAt,
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
			`INSERT IGNORE INTO t_tx_btc_token (
    id,
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc_token (
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
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

// SQLGetTTxBtcToken 根据id查询
func SQLGetTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTTxBtcToken, error) {
	var row DBTTxBtcToken
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
FROM
	t_tx_btc_token
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

// SQLGetTTxBtcTokenCol 根据id查询
func SQLGetTTxBtcTokenCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTTxBtcToken, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_token
WHERE
	id=:id`)

	var row DBTTxBtcToken
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

// SQLSelectTTxBtcToken 根据ids获取
func SQLSelectTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTTxBtcToken, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTTxBtcToken
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    token_index,
    token_symbol,
    block_hash,
    tx_id,
    from_address,
    to_address,
    value,
    blocktime,
    create_at,
    handle_status,
    handle_msg,
    handle_at
FROM
	t_tx_btc_token
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

// SQLSelectTTxBtcTokenCol 根据ids获取
func SQLSelectTTxBtcTokenCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTTxBtcToken, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_token
WHERE
	id IN (:ids)`)

	var rows []*DBTTxBtcToken
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

// SQLUpdateTTxBtcToken 更新
func SQLUpdateTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtcToken) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_btc_token
SET
    token_index=:token_index,
    token_symbol=:token_symbol,
    block_hash=:block_hash,
    tx_id=:tx_id,
    from_address=:from_address,
    to_address=:to_address,
    value=:value,
    blocktime=:blocktime,
    create_at=:create_at,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_at=:handle_at
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"token_index":   row.TokenIndex,
			"token_symbol":  row.TokenSymbol,
			"block_hash":    row.BlockHash,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"value":         row.Value,
			"blocktime":     row.Blocktime,
			"create_at":     row.CreateAt,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_at":     row.HandleAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTTxBtcToken 删除
func SQLDeleteTTxBtcToken(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_btc_token
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

// SQLCreateTTxBtcUxto 创建
func SQLCreateTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtcUxto) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc_uxto (
    id,
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :uxto_type,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :vout_script,
    :create_time,
    :spend_tx_id,
    :spend_n,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"uxto_type":     row.UxtoType,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"vout_script":   row.VoutScript,
				"create_time":   row.CreateTime,
				"spend_tx_id":   row.SpendTxID,
				"spend_n":       row.SpendN,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc_uxto (
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :uxto_type,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :vout_script,
    :create_time,
    :spend_tx_id,
    :spend_n,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"uxto_type":     row.UxtoType,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"vout_script":   row.VoutScript,
				"create_time":   row.CreateTime,
				"spend_tx_id":   row.SpendTxID,
				"spend_n":       row.SpendN,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTTxBtcUxto 创建
func SQLCreateIgnoreTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtcUxto) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc_uxto (
    id,
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :uxto_type,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :vout_script,
    :create_time,
    :spend_tx_id,
    :spend_n,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"uxto_type":     row.UxtoType,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"vout_script":   row.VoutScript,
				"create_time":   row.CreateTime,
				"spend_tx_id":   row.SpendTxID,
				"spend_n":       row.SpendN,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc_uxto (
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :uxto_type,
    :block_hash,
    :tx_id,
    :vout_n,
    :vout_address,
    :vout_value,
    :vout_script,
    :create_time,
    :spend_tx_id,
    :spend_n,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"uxto_type":     row.UxtoType,
				"block_hash":    row.BlockHash,
				"tx_id":         row.TxID,
				"vout_n":        row.VoutN,
				"vout_address":  row.VoutAddress,
				"vout_value":    row.VoutValue,
				"vout_script":   row.VoutScript,
				"create_time":   row.CreateTime,
				"spend_tx_id":   row.SpendTxID,
				"spend_n":       row.SpendN,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTTxBtcUxto 创建多个
func SQLCreateManyTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxBtcUxto) (int64, error) {
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
					row.UxtoType,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.VoutScript,
					row.CreateTime,
					row.SpendTxID,
					row.SpendN,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.UxtoType,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.VoutScript,
					row.CreateTime,
					row.SpendTxID,
					row.SpendN,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT INTO t_tx_btc_uxto (
    id,
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_tx_btc_uxto (
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
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

// SQLCreateIgnoreManyTTxBtcUxto 创建多个
func SQLCreateIgnoreManyTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxBtcUxto) (int64, error) {
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
					row.UxtoType,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.VoutScript,
					row.CreateTime,
					row.SpendTxID,
					row.SpendN,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.UxtoType,
					row.BlockHash,
					row.TxID,
					row.VoutN,
					row.VoutAddress,
					row.VoutValue,
					row.VoutScript,
					row.CreateTime,
					row.SpendTxID,
					row.SpendN,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT IGNORE INTO t_tx_btc_uxto (
    id,
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_btc_uxto (
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
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

// SQLGetTTxBtcUxto 根据id查询
func SQLGetTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTTxBtcUxto, error) {
	var row DBTTxBtcUxto
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_tx_btc_uxto
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

// SQLGetTTxBtcUxtoCol 根据id查询
func SQLGetTTxBtcUxtoCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTTxBtcUxto, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_uxto
WHERE
	id=:id`)

	var row DBTTxBtcUxto
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

// SQLSelectTTxBtcUxto 根据ids获取
func SQLSelectTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTTxBtcUxto, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTTxBtcUxto
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    uxto_type,
    block_hash,
    tx_id,
    vout_n,
    vout_address,
    vout_value,
    vout_script,
    create_time,
    spend_tx_id,
    spend_n,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_tx_btc_uxto
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

// SQLSelectTTxBtcUxtoCol 根据ids获取
func SQLSelectTTxBtcUxtoCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTTxBtcUxto, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_uxto
WHERE
	id IN (:ids)`)

	var rows []*DBTTxBtcUxto
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

// SQLUpdateTTxBtcUxto 更新
func SQLUpdateTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxBtcUxto) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_btc_uxto
SET
    uxto_type=:uxto_type,
    block_hash=:block_hash,
    tx_id=:tx_id,
    vout_n=:vout_n,
    vout_address=:vout_address,
    vout_value=:vout_value,
    vout_script=:vout_script,
    create_time=:create_time,
    spend_tx_id=:spend_tx_id,
    spend_n=:spend_n,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"uxto_type":     row.UxtoType,
			"block_hash":    row.BlockHash,
			"tx_id":         row.TxID,
			"vout_n":        row.VoutN,
			"vout_address":  row.VoutAddress,
			"vout_value":    row.VoutValue,
			"vout_script":   row.VoutScript,
			"create_time":   row.CreateTime,
			"spend_tx_id":   row.SpendTxID,
			"spend_n":       row.SpendN,
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

// SQLDeleteTTxBtcUxto 删除
func SQLDeleteTTxBtcUxto(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_btc_uxto
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

// SQLCreateTTxErc20 创建
func SQLCreateTTxErc20(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxErc20) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_tx_erc20 (
    id,
    token_id,
    product_id,
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
    :token_id,
    :product_id,
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
				"token_id":      row.TokenID,
				"product_id":    row.ProductID,
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
			`INSERT INTO t_tx_erc20 (
    token_id,
    product_id,
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
    :token_id,
    :product_id,
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
				"token_id":      row.TokenID,
				"product_id":    row.ProductID,
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

// SQLCreateIgnoreTTxErc20 创建
func SQLCreateIgnoreTTxErc20(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxErc20) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_tx_erc20 (
    id,
    token_id,
    product_id,
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
    :token_id,
    :product_id,
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
				"token_id":      row.TokenID,
				"product_id":    row.ProductID,
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
			`INSERT IGNORE INTO t_tx_erc20 (
    token_id,
    product_id,
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
    :token_id,
    :product_id,
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
				"token_id":      row.TokenID,
				"product_id":    row.ProductID,
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

// SQLCreateManyTTxErc20 创建多个
func SQLCreateManyTTxErc20(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxErc20) (int64, error) {
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
					row.TokenID,
					row.ProductID,
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
					row.TokenID,
					row.ProductID,
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
			`INSERT INTO t_tx_erc20 (
    id,
    token_id,
    product_id,
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
			`INSERT INTO t_tx_erc20 (
    token_id,
    product_id,
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

// SQLCreateIgnoreManyTTxErc20 创建多个
func SQLCreateIgnoreManyTTxErc20(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTTxErc20) (int64, error) {
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
					row.TokenID,
					row.ProductID,
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
					row.TokenID,
					row.ProductID,
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
			`INSERT IGNORE INTO t_tx_erc20 (
    id,
    token_id,
    product_id,
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
			`INSERT IGNORE INTO t_tx_erc20 (
    token_id,
    product_id,
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

// SQLGetTTxErc20 根据id查询
func SQLGetTTxErc20(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTTxErc20, error) {
	var row DBTTxErc20
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    token_id,
    product_id,
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
	t_tx_erc20
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

// SQLGetTTxErc20Col 根据id查询
func SQLGetTTxErc20Col(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTTxErc20, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
WHERE
	id=:id`)

	var row DBTTxErc20
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

// SQLSelectTTxErc20 根据ids获取
func SQLSelectTTxErc20(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTTxErc20, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTTxErc20
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    token_id,
    product_id,
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
	t_tx_erc20
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

// SQLSelectTTxErc20Col 根据ids获取
func SQLSelectTTxErc20Col(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTTxErc20, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
WHERE
	id IN (:ids)`)

	var rows []*DBTTxErc20
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

// SQLUpdateTTxErc20 更新
func SQLUpdateTTxErc20(ctx context.Context, tx hcommon.DbExeAble, row *DBTTxErc20) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_erc20
SET
    token_id=:token_id,
    product_id=:product_id,
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
			"token_id":      row.TokenID,
			"product_id":    row.ProductID,
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

// SQLDeleteTTxErc20 删除
func SQLDeleteTTxErc20(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_erc20
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

// SQLCreateTWithdraw 创建
func SQLCreateTWithdraw(ctx context.Context, tx hcommon.DbExeAble, row *DBTWithdraw) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_withdraw (
    id,
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :product_id,
    :out_serial,
    :to_address,
    :symbol,
    :balance_real,
    :tx_hash,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"product_id":    row.ProductID,
				"out_serial":    row.OutSerial,
				"to_address":    row.ToAddress,
				"symbol":        row.Symbol,
				"balance_real":  row.BalanceReal,
				"tx_hash":       row.TxHash,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT INTO t_withdraw (
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :product_id,
    :out_serial,
    :to_address,
    :symbol,
    :balance_real,
    :tx_hash,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"product_id":    row.ProductID,
				"out_serial":    row.OutSerial,
				"to_address":    row.ToAddress,
				"symbol":        row.Symbol,
				"balance_real":  row.BalanceReal,
				"tx_hash":       row.TxHash,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateIgnoreTWithdraw 创建
func SQLCreateIgnoreTWithdraw(ctx context.Context, tx hcommon.DbExeAble, row *DBTWithdraw) (int64, error) {
	var lastID int64
	var err error
	if row.ID > 0 {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_withdraw (
    id,
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :id,
    :product_id,
    :out_serial,
    :to_address,
    :symbol,
    :balance_real,
    :tx_hash,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"id":            row.ID,
				"product_id":    row.ProductID,
				"out_serial":    row.OutSerial,
				"to_address":    row.ToAddress,
				"symbol":        row.Symbol,
				"balance_real":  row.BalanceReal,
				"tx_hash":       row.TxHash,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	} else {
		lastID, err = hcommon.DbExecuteLastIDNamedContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_withdraw (
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES (
    :product_id,
    :out_serial,
    :to_address,
    :symbol,
    :balance_real,
    :tx_hash,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`,
			gin.H{
				"product_id":    row.ProductID,
				"out_serial":    row.OutSerial,
				"to_address":    row.ToAddress,
				"symbol":        row.Symbol,
				"balance_real":  row.BalanceReal,
				"tx_hash":       row.TxHash,
				"create_time":   row.CreateTime,
				"handle_status": row.HandleStatus,
				"handle_msg":    row.HandleMsg,
				"handle_time":   row.HandleTime,
			},
		)
	}
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTWithdraw 创建多个
func SQLCreateManyTWithdraw(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTWithdraw) (int64, error) {
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
					row.ProductID,
					row.OutSerial,
					row.ToAddress,
					row.Symbol,
					row.BalanceReal,
					row.TxHash,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ProductID,
					row.OutSerial,
					row.ToAddress,
					row.Symbol,
					row.BalanceReal,
					row.TxHash,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT INTO t_withdraw (
    id,
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT INTO t_withdraw (
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLCreateIgnoreManyTWithdraw 创建多个
func SQLCreateIgnoreManyTWithdraw(ctx context.Context, tx hcommon.DbExeAble, rows []*DBTWithdraw) (int64, error) {
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
					row.ProductID,
					row.OutSerial,
					row.ToAddress,
					row.Symbol,
					row.BalanceReal,
					row.TxHash,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ProductID,
					row.OutSerial,
					row.ToAddress,
					row.Symbol,
					row.BalanceReal,
					row.TxHash,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleTime,
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
			`INSERT IGNORE INTO t_withdraw (
    id,
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`,
			len(rows),
			args...,
		)
	} else {
		count, err = hcommon.DbExecuteCountManyContent(
			ctx,
			tx,
			`INSERT IGNORE INTO t_withdraw (
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
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

// SQLGetTWithdraw 根据id查询
func SQLGetTWithdraw(ctx context.Context, tx hcommon.DbExeAble, id int64) (*DBTWithdraw, error) {
	var row DBTWithdraw
	ok, err := hcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		`SELECT
    id,
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_withdraw
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

// SQLGetTWithdrawCol 根据id查询
func SQLGetTWithdrawCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, id int64) (*DBTWithdraw, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	id=:id`)

	var row DBTWithdraw
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

// SQLSelectTWithdraw 根据ids获取
func SQLSelectTWithdraw(ctx context.Context, tx hcommon.DbExeAble, ids []int64) ([]*DBTWithdraw, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var rows []*DBTWithdraw
	err := hcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		`SELECT
    id,
    product_id,
    out_serial,
    to_address,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
FROM
	t_withdraw
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

// SQLSelectTWithdrawCol 根据ids获取
func SQLSelectTWithdrawCol(ctx context.Context, tx hcommon.DbExeAble, cols []string, ids []int64) ([]*DBTWithdraw, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	id IN (:ids)`)

	var rows []*DBTWithdraw
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

// SQLUpdateTWithdraw 更新
func SQLUpdateTWithdraw(ctx context.Context, tx hcommon.DbExeAble, row *DBTWithdraw) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
SET
    product_id=:product_id,
    out_serial=:out_serial,
    to_address=:to_address,
    symbol=:symbol,
    balance_real=:balance_real,
    tx_hash=:tx_hash,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		gin.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"out_serial":    row.OutSerial,
			"to_address":    row.ToAddress,
			"symbol":        row.Symbol,
			"balance_real":  row.BalanceReal,
			"tx_hash":       row.TxHash,
			"create_time":   row.CreateTime,
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

// SQLDeleteTWithdraw 删除
func SQLDeleteTWithdraw(ctx context.Context, tx hcommon.DbExeAble, id int64) (int64, error) {
	count, err := hcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_withdraw
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
