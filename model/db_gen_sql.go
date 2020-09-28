package model

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/moremorefun/mcommon"
)

// SQLCreateTAddressKey 创建
func SQLCreateTAddressKey(ctx context.Context, tx mcommon.DbExeAble, row *DBTAddressKey, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_address_key ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       symbol,
       address,
       pwd,
       use_tag
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :symbol,
    :address,
    :pwd,
    :use_tag
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateTAddressKeyDuplicate 创建更新
func SQLCreateTAddressKeyDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAddressKey, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_address_key ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       symbol,
       address,
       pwd,
       use_tag
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :symbol,
    :address,
    :pwd,
    :use_tag
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateManyTAddressKey 创建多个
func SQLCreateManyTAddressKey(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAddressKey, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_address_key ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    symbol,
    address,
    pwd,
    use_tag
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAddressKeyDuplicate 创建多个
func SQLCreateManyTAddressKeyDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAddressKey, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_address_key ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    symbol,
    address,
    pwd,
    use_tag
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAddressKeyCol 根据id查询
func SQLGetTAddressKeyCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAddressKey, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
WHERE
	id=:id`)

	var row DBTAddressKey
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAddressKeyColKV 根据id查询
func SQLGetTAddressKeyColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAddressKey, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAddressKey
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAddressKeyCol 根据ids获取
func SQLSelectTAddressKeyCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAddressKey, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAddressKey
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAddressKeyColKV 根据ids获取
func SQLSelectTAddressKeyColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAddressKey, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_address_key
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAddressKey
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAddressKey 更新
func SQLUpdateTAddressKey(ctx context.Context, tx mcommon.DbExeAble, row *DBTAddressKey) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
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
func SQLDeleteTAddressKey(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_address_key
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppConfigInt 创建
func SQLCreateTAppConfigInt(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigInt, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_int ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTAppConfigIntDuplicate 创建更新
func SQLCreateTAppConfigIntDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigInt, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_int ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppConfigInt 创建多个
func SQLCreateManyTAppConfigInt(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigInt, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_int ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAppConfigIntDuplicate 创建多个
func SQLCreateManyTAppConfigIntDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigInt, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_int ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppConfigIntCol 根据id查询
func SQLGetTAppConfigIntCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigInt, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_int
WHERE
	id=:id`)

	var row DBTAppConfigInt
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAppConfigIntColKV 根据id查询
func SQLGetTAppConfigIntColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAppConfigInt, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_int
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAppConfigInt
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAppConfigIntCol 根据ids获取
func SQLSelectTAppConfigIntCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAppConfigInt, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAppConfigInt
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppConfigIntColKV 根据ids获取
func SQLSelectTAppConfigIntColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAppConfigInt, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_int
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAppConfigInt
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppConfigInt 更新
func SQLUpdateTAppConfigInt(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigInt) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_int
SET
    k=:k,
    v=:v
WHERE
	id=:id`,
		mcommon.H{
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
func SQLDeleteTAppConfigInt(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_int
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppConfigStr 创建
func SQLCreateTAppConfigStr(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigStr, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_str ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTAppConfigStrDuplicate 创建更新
func SQLCreateTAppConfigStrDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigStr, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_str ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppConfigStr 创建多个
func SQLCreateManyTAppConfigStr(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigStr, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_str ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAppConfigStrDuplicate 创建多个
func SQLCreateManyTAppConfigStrDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigStr, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_str ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppConfigStrCol 根据id查询
func SQLGetTAppConfigStrCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigStr, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_str
WHERE
	id=:id`)

	var row DBTAppConfigStr
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAppConfigStrColKV 根据id查询
func SQLGetTAppConfigStrColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAppConfigStr, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_str
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAppConfigStr
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAppConfigStrCol 根据ids获取
func SQLSelectTAppConfigStrCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAppConfigStr, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAppConfigStr
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppConfigStrColKV 根据ids获取
func SQLSelectTAppConfigStrColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAppConfigStr, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_str
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAppConfigStr
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppConfigStr 更新
func SQLUpdateTAppConfigStr(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigStr) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_str
SET
    k=:k,
    v=:v
WHERE
	id=:id`,
		mcommon.H{
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
func SQLDeleteTAppConfigStr(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_str
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppConfigToken 创建
func SQLCreateTAppConfigToken(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigToken, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_token ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       token_address,
       token_decimals,
       token_symbol,
       cold_address,
       hot_address,
       org_min_balance,
       create_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :token_address,
    :token_decimals,
    :token_symbol,
    :cold_address,
    :hot_address,
    :org_min_balance,
    :create_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateTAppConfigTokenDuplicate 创建更新
func SQLCreateTAppConfigTokenDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigToken, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_token ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       token_address,
       token_decimals,
       token_symbol,
       cold_address,
       hot_address,
       org_min_balance,
       create_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :token_address,
    :token_decimals,
    :token_symbol,
    :cold_address,
    :hot_address,
    :org_min_balance,
    :create_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateManyTAppConfigToken 创建多个
func SQLCreateManyTAppConfigToken(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigToken, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_token ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAppConfigTokenDuplicate 创建多个
func SQLCreateManyTAppConfigTokenDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigToken, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_token ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    token_address,
    token_decimals,
    token_symbol,
    cold_address,
    hot_address,
    org_min_balance,
    create_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppConfigTokenCol 根据id查询
func SQLGetTAppConfigTokenCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigToken, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token
WHERE
	id=:id`)

	var row DBTAppConfigToken
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAppConfigTokenColKV 根据id查询
func SQLGetTAppConfigTokenColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAppConfigToken, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAppConfigToken
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAppConfigTokenCol 根据ids获取
func SQLSelectTAppConfigTokenCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAppConfigToken, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAppConfigToken
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppConfigTokenColKV 根据ids获取
func SQLSelectTAppConfigTokenColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAppConfigToken, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAppConfigToken
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppConfigToken 更新
func SQLUpdateTAppConfigToken(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigToken) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
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
func SQLDeleteTAppConfigToken(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_token
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppConfigTokenBtc 创建
func SQLCreateTAppConfigTokenBtc(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigTokenBtc, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_token_btc ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       token_index,
       token_symbol,
       cold_address,
       hot_address,
       fee_address,
       tx_org_min_balance,
       create_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :token_index,
    :token_symbol,
    :cold_address,
    :hot_address,
    :fee_address,
    :tx_org_min_balance,
    :create_at
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":                 row.ID,
			"token_index":        row.TokenIndex,
			"token_symbol":       row.TokenSymbol,
			"cold_address":       row.ColdAddress,
			"hot_address":        row.HotAddress,
			"fee_address":        row.FeeAddress,
			"tx_org_min_balance": row.TxOrgMinBalance,
			"create_at":          row.CreateAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTAppConfigTokenBtcDuplicate 创建更新
func SQLCreateTAppConfigTokenBtcDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigTokenBtc, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_token_btc ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       token_index,
       token_symbol,
       cold_address,
       hot_address,
       fee_address,
       tx_org_min_balance,
       create_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :token_index,
    :token_symbol,
    :cold_address,
    :hot_address,
    :fee_address,
    :tx_org_min_balance,
    :create_at
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":                 row.ID,
			"token_index":        row.TokenIndex,
			"token_symbol":       row.TokenSymbol,
			"cold_address":       row.ColdAddress,
			"hot_address":        row.HotAddress,
			"fee_address":        row.FeeAddress,
			"tx_org_min_balance": row.TxOrgMinBalance,
			"create_at":          row.CreateAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppConfigTokenBtc 创建多个
func SQLCreateManyTAppConfigTokenBtc(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigTokenBtc, isIgnore bool) (int64, error) {
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
					row.FeeAddress,
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
					row.FeeAddress,
					row.TxOrgMinBalance,
					row.CreateAt,
				},
			)
		}
	}
	var count int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_config_token_btc ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    fee_address,
    tx_org_min_balance,
    create_at
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAppConfigTokenBtcDuplicate 创建多个
func SQLCreateManyTAppConfigTokenBtcDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppConfigTokenBtc, updates []string) (int64, error) {
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
					row.FeeAddress,
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
					row.FeeAddress,
					row.TxOrgMinBalance,
					row.CreateAt,
				},
			)
		}
	}
	var count int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_config_token_btc ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    token_index,
    token_symbol,
    cold_address,
    hot_address,
    fee_address,
    tx_org_min_balance,
    create_at
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppConfigTokenBtcCol 根据id查询
func SQLGetTAppConfigTokenBtcCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAppConfigTokenBtc, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token_btc
WHERE
	id=:id`)

	var row DBTAppConfigTokenBtc
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAppConfigTokenBtcColKV 根据id查询
func SQLGetTAppConfigTokenBtcColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAppConfigTokenBtc, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token_btc
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAppConfigTokenBtc
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAppConfigTokenBtcCol 根据ids获取
func SQLSelectTAppConfigTokenBtcCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAppConfigTokenBtc, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAppConfigTokenBtc
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppConfigTokenBtcColKV 根据ids获取
func SQLSelectTAppConfigTokenBtcColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAppConfigTokenBtc, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_config_token_btc
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAppConfigTokenBtc
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppConfigTokenBtc 更新
func SQLUpdateTAppConfigTokenBtc(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppConfigTokenBtc) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_config_token_btc
SET
    token_index=:token_index,
    token_symbol=:token_symbol,
    cold_address=:cold_address,
    hot_address=:hot_address,
    fee_address=:fee_address,
    tx_org_min_balance=:tx_org_min_balance,
    create_at=:create_at
WHERE
	id=:id`,
		mcommon.H{
			"id":                 row.ID,
			"token_index":        row.TokenIndex,
			"token_symbol":       row.TokenSymbol,
			"cold_address":       row.ColdAddress,
			"hot_address":        row.HotAddress,
			"fee_address":        row.FeeAddress,
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
func SQLDeleteTAppConfigTokenBtc(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_config_token_btc
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppLock 创建
func SQLCreateTAppLock(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppLock, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_lock ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v,
       create_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v,
    :create_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":          row.ID,
			"k":           row.K,
			"v":           row.V,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTAppLockDuplicate 创建更新
func SQLCreateTAppLockDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppLock, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_lock ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v,
       create_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v,
    :create_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":          row.ID,
			"k":           row.K,
			"v":           row.V,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppLock 创建多个
func SQLCreateManyTAppLock(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppLock, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_lock ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v,
    create_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAppLockDuplicate 创建多个
func SQLCreateManyTAppLockDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppLock, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_lock ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v,
    create_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppLockCol 根据id查询
func SQLGetTAppLockCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAppLock, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
WHERE
	id=:id`)

	var row DBTAppLock
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAppLockColKV 根据id查询
func SQLGetTAppLockColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAppLock, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAppLock
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAppLockCol 根据ids获取
func SQLSelectTAppLockCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAppLock, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAppLock
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppLockColKV 根据ids获取
func SQLSelectTAppLockColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAppLock, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_lock
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAppLock
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppLock 更新
func SQLUpdateTAppLock(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppLock) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
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
func SQLDeleteTAppLock(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_lock
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTAppStatusInt 创建
func SQLCreateTAppStatusInt(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppStatusInt, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_status_int ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTAppStatusIntDuplicate 创建更新
func SQLCreateTAppStatusIntDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppStatusInt, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_status_int ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       k,
       v
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :k,
    :v
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id": row.ID,
			"k":  row.K,
			"v":  row.V,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTAppStatusInt 创建多个
func SQLCreateManyTAppStatusInt(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppStatusInt, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_app_status_int ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTAppStatusIntDuplicate 创建多个
func SQLCreateManyTAppStatusIntDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTAppStatusInt, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_app_status_int ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    k,
    v
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTAppStatusIntCol 根据id查询
func SQLGetTAppStatusIntCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTAppStatusInt, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_status_int
WHERE
	id=:id`)

	var row DBTAppStatusInt
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTAppStatusIntColKV 根据id查询
func SQLGetTAppStatusIntColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTAppStatusInt, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_status_int
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTAppStatusInt
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTAppStatusIntCol 根据ids获取
func SQLSelectTAppStatusIntCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTAppStatusInt, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTAppStatusInt
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTAppStatusIntColKV 根据ids获取
func SQLSelectTAppStatusIntColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTAppStatusInt, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_app_status_int
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTAppStatusInt
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTAppStatusInt 更新
func SQLUpdateTAppStatusInt(ctx context.Context, tx mcommon.DbExeAble, row *DBTAppStatusInt) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_app_status_int
SET
    k=:k,
    v=:v
WHERE
	id=:id`,
		mcommon.H{
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
func SQLDeleteTAppStatusInt(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_app_status_int
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTProduct 创建
func SQLCreateTProduct(ctx context.Context, tx mcommon.DbExeAble, row *DBTProduct, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_product ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       app_name,
       app_sk,
       cb_url,
       whitelist_ip
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :app_name,
    :app_sk,
    :cb_url,
    :whitelist_ip
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateTProductDuplicate 创建更新
func SQLCreateTProductDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTProduct, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_product ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       app_name,
       app_sk,
       cb_url,
       whitelist_ip
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :app_name,
    :app_sk,
    :cb_url,
    :whitelist_ip
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateManyTProduct 创建多个
func SQLCreateManyTProduct(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTProduct, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_product ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTProductDuplicate 创建多个
func SQLCreateManyTProductDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTProduct, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_product ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    app_name,
    app_sk,
    cb_url,
    whitelist_ip
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTProductCol 根据id查询
func SQLGetTProductCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTProduct, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product
WHERE
	id=:id`)

	var row DBTProduct
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTProductColKV 根据id查询
func SQLGetTProductColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTProduct, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTProduct
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTProductCol 根据ids获取
func SQLSelectTProductCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTProduct, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTProduct
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTProductColKV 根据ids获取
func SQLSelectTProductColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTProduct, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTProduct
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTProduct 更新
func SQLUpdateTProduct(ctx context.Context, tx mcommon.DbExeAble, row *DBTProduct) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
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
func SQLDeleteTProduct(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_product
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTProductNonce 创建
func SQLCreateTProductNonce(ctx context.Context, tx mcommon.DbExeAble, row *DBTProductNonce, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_product_nonce ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       c,
       create_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :c,
    :create_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":          row.ID,
			"c":           row.C,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTProductNonceDuplicate 创建更新
func SQLCreateTProductNonceDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTProductNonce, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_product_nonce ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       c,
       create_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :c,
    :create_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":          row.ID,
			"c":           row.C,
			"create_time": row.CreateTime,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTProductNonce 创建多个
func SQLCreateManyTProductNonce(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTProductNonce, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_product_nonce ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    c,
    create_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTProductNonceDuplicate 创建多个
func SQLCreateManyTProductNonceDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTProductNonce, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_product_nonce ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    c,
    create_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTProductNonceCol 根据id查询
func SQLGetTProductNonceCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTProductNonce, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_nonce
WHERE
	id=:id`)

	var row DBTProductNonce
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTProductNonceColKV 根据id查询
func SQLGetTProductNonceColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTProductNonce, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_nonce
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTProductNonce
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTProductNonceCol 根据ids获取
func SQLSelectTProductNonceCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTProductNonce, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTProductNonce
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTProductNonceColKV 根据ids获取
func SQLSelectTProductNonceColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTProductNonce, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_nonce
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTProductNonce
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTProductNonce 更新
func SQLUpdateTProductNonce(ctx context.Context, tx mcommon.DbExeAble, row *DBTProductNonce) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_product_nonce
SET
    c=:c,
    create_time=:create_time
WHERE
	id=:id`,
		mcommon.H{
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
func SQLDeleteTProductNonce(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_product_nonce
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTProductNotify 创建
func SQLCreateTProductNotify(ctx context.Context, tx mcommon.DbExeAble, row *DBTProductNotify, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_product_notify ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       nonce,
       product_id,
       item_type,
       item_id,
       notify_type,
       token_symbol,
       url,
       msg,
       handle_status,
       handle_msg,
       create_time,
       update_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :nonce,
    :product_id,
    :item_type,
    :item_id,
    :notify_type,
    :token_symbol,
    :url,
    :msg,
    :handle_status,
    :handle_msg,
    :create_time,
    :update_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"nonce":         row.Nonce,
			"product_id":    row.ProductID,
			"item_type":     row.ItemType,
			"item_id":       row.ItemID,
			"notify_type":   row.NotifyType,
			"token_symbol":  row.TokenSymbol,
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
	return lastID, nil
}

// SQLCreateTProductNotifyDuplicate 创建更新
func SQLCreateTProductNotifyDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTProductNotify, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_product_notify ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       nonce,
       product_id,
       item_type,
       item_id,
       notify_type,
       token_symbol,
       url,
       msg,
       handle_status,
       handle_msg,
       create_time,
       update_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :nonce,
    :product_id,
    :item_type,
    :item_id,
    :notify_type,
    :token_symbol,
    :url,
    :msg,
    :handle_status,
    :handle_msg,
    :create_time,
    :update_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"nonce":         row.Nonce,
			"product_id":    row.ProductID,
			"item_type":     row.ItemType,
			"item_id":       row.ItemID,
			"notify_type":   row.NotifyType,
			"token_symbol":  row.TokenSymbol,
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
	return lastID, nil
}

// SQLCreateManyTProductNotify 创建多个
func SQLCreateManyTProductNotify(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTProductNotify, isIgnore bool) (int64, error) {
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
					row.TokenSymbol,
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
					row.TokenSymbol,
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_product_notify ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    token_symbol,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTProductNotifyDuplicate 创建多个
func SQLCreateManyTProductNotifyDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTProductNotify, updates []string) (int64, error) {
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
					row.TokenSymbol,
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
					row.TokenSymbol,
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_product_notify ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    nonce,
    product_id,
    item_type,
    item_id,
    notify_type,
    token_symbol,
    url,
    msg,
    handle_status,
    handle_msg,
    create_time,
    update_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTProductNotifyCol 根据id查询
func SQLGetTProductNotifyCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTProductNotify, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
WHERE
	id=:id`)

	var row DBTProductNotify
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTProductNotifyColKV 根据id查询
func SQLGetTProductNotifyColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTProductNotify, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTProductNotify
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTProductNotifyCol 根据ids获取
func SQLSelectTProductNotifyCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTProductNotify, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTProductNotify
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTProductNotifyColKV 根据ids获取
func SQLSelectTProductNotifyColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTProductNotify, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_product_notify
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTProductNotify
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTProductNotify 更新
func SQLUpdateTProductNotify(ctx context.Context, tx mcommon.DbExeAble, row *DBTProductNotify) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
    token_symbol=:token_symbol,
    url=:url,
    msg=:msg,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    create_time=:create_time,
    update_time=:update_time
WHERE
	id=:id`,
		mcommon.H{
			"id":            row.ID,
			"nonce":         row.Nonce,
			"product_id":    row.ProductID,
			"item_type":     row.ItemType,
			"item_id":       row.ItemID,
			"notify_type":   row.NotifyType,
			"token_symbol":  row.TokenSymbol,
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
func SQLDeleteTProductNotify(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_product_notify
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTSend 创建
func SQLCreateTSend(ctx context.Context, tx mcommon.DbExeAble, row *DBTSend, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       related_type,
       related_id,
       token_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       gas,
       gas_price,
       nonce,
       hex,
       create_time,
       handle_status,
       handle_msg,
       handle_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :gas,
    :gas_price,
    :nonce,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateTSendDuplicate 创建更新
func SQLCreateTSendDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTSend, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_send ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       related_type,
       related_id,
       token_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       gas,
       gas_price,
       nonce,
       hex,
       create_time,
       handle_status,
       handle_msg,
       handle_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :gas,
    :gas_price,
    :nonce,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateManyTSend 创建多个
func SQLCreateManyTSend(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTSend, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
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
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTSendDuplicate 创建多个
func SQLCreateManyTSendDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTSend, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_send ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
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
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTSendCol 根据id查询
func SQLGetTSendCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTSend, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send
WHERE
	id=:id`)

	var row DBTSend
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTSendColKV 根据id查询
func SQLGetTSendColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTSend, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTSend
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTSendCol 根据ids获取
func SQLSelectTSendCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTSend, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTSend
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTSendColKV 根据ids获取
func SQLSelectTSendColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTSend, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTSend
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTSend 更新
func SQLUpdateTSend(ctx context.Context, tx mcommon.DbExeAble, row *DBTSend) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
func SQLDeleteTSend(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_send
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTSendBtc 创建
func SQLCreateTSendBtc(ctx context.Context, tx mcommon.DbExeAble, row *DBTSendBtc, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send_btc ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       related_type,
       related_id,
       token_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       gas,
       gas_price,
       hex,
       create_time,
       handle_status,
       handle_msg,
       handle_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :gas,
    :gas_price,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateTSendBtcDuplicate 创建更新
func SQLCreateTSendBtcDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTSendBtc, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_send_btc ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       related_type,
       related_id,
       token_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       gas,
       gas_price,
       hex,
       create_time,
       handle_status,
       handle_msg,
       handle_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :related_type,
    :related_id,
    :token_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :gas,
    :gas_price,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateManyTSendBtc 创建多个
func SQLCreateManyTSendBtc(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTSendBtc, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send_btc ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTSendBtcDuplicate 创建多个
func SQLCreateManyTSendBtcDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTSendBtc, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_send_btc ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    related_type,
    related_id,
    token_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    gas,
    gas_price,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTSendBtcCol 根据id查询
func SQLGetTSendBtcCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTSendBtc, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_btc
WHERE
	id=:id`)

	var row DBTSendBtc
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTSendBtcColKV 根据id查询
func SQLGetTSendBtcColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTSendBtc, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_btc
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTSendBtc
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTSendBtcCol 根据ids获取
func SQLSelectTSendBtcCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTSendBtc, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTSendBtc
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTSendBtcColKV 根据ids获取
func SQLSelectTSendBtcColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTSendBtc, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_btc
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTSendBtc
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTSendBtc 更新
func SQLUpdateTSendBtc(ctx context.Context, tx mcommon.DbExeAble, row *DBTSendBtc) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
			"id":            row.ID,
			"related_type":  row.RelatedType,
			"related_id":    row.RelatedID,
			"token_id":      row.TokenID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
func SQLDeleteTSendBtc(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_send_btc
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTSendEos 创建
func SQLCreateTSendEos(ctx context.Context, tx mcommon.DbExeAble, row *DBTSendEos, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send_eos ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       withdraw_id,
       tx_hash,
       log_index,
       from_address,
       to_address,
       memo,
       balance_real,
       hex,
       create_time,
       handle_status,
       handle_msg,
       handle_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :withdraw_id,
    :tx_hash,
    :log_index,
    :from_address,
    :to_address,
    :memo,
    :balance_real,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_at
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"withdraw_id":   row.WithdrawID,
			"tx_hash":       row.TxHash,
			"log_index":     row.LogIndex,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
			"balance_real":  row.BalanceReal,
			"hex":           row.Hex,
			"create_time":   row.CreateTime,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_at":     row.HandleAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTSendEosDuplicate 创建更新
func SQLCreateTSendEosDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTSendEos, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_send_eos ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       withdraw_id,
       tx_hash,
       log_index,
       from_address,
       to_address,
       memo,
       balance_real,
       hex,
       create_time,
       handle_status,
       handle_msg,
       handle_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :withdraw_id,
    :tx_hash,
    :log_index,
    :from_address,
    :to_address,
    :memo,
    :balance_real,
    :hex,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_at
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"withdraw_id":   row.WithdrawID,
			"tx_hash":       row.TxHash,
			"log_index":     row.LogIndex,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
			"balance_real":  row.BalanceReal,
			"hex":           row.Hex,
			"create_time":   row.CreateTime,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_at":     row.HandleAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTSendEos 创建多个
func SQLCreateManyTSendEos(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTSendEos, isIgnore bool) (int64, error) {
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
					row.WithdrawID,
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
					row.Hex,
					row.CreateTime,
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
					row.WithdrawID,
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleAt,
				},
			)
		}
	}
	var count int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_send_eos ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    withdraw_id,
    tx_hash,
    log_index,
    from_address,
    to_address,
    memo,
    balance_real,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_at
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTSendEosDuplicate 创建多个
func SQLCreateManyTSendEosDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTSendEos, updates []string) (int64, error) {
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
					row.WithdrawID,
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
					row.Hex,
					row.CreateTime,
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
					row.WithdrawID,
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
					row.Hex,
					row.CreateTime,
					row.HandleStatus,
					row.HandleMsg,
					row.HandleAt,
				},
			)
		}
	}
	var count int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_send_eos ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    withdraw_id,
    tx_hash,
    log_index,
    from_address,
    to_address,
    memo,
    balance_real,
    hex,
    create_time,
    handle_status,
    handle_msg,
    handle_at
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTSendEosCol 根据id查询
func SQLGetTSendEosCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTSendEos, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_eos
WHERE
	id=:id`)

	var row DBTSendEos
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTSendEosColKV 根据id查询
func SQLGetTSendEosColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTSendEos, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_eos
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTSendEos
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTSendEosCol 根据ids获取
func SQLSelectTSendEosCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTSendEos, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_eos
WHERE
	id IN (:ids)`)
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTSendEos
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTSendEosColKV 根据ids获取
func SQLSelectTSendEosColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTSendEos, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_send_eos
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTSendEos
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTSendEos 更新
func SQLUpdateTSendEos(ctx context.Context, tx mcommon.DbExeAble, row *DBTSendEos) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_send_eos
SET
    withdraw_id=:withdraw_id,
    tx_hash=:tx_hash,
    log_index=:log_index,
    from_address=:from_address,
    to_address=:to_address,
    memo=:memo,
    balance_real=:balance_real,
    hex=:hex,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_at=:handle_at
WHERE
	id=:id`,
		mcommon.H{
			"id":            row.ID,
			"withdraw_id":   row.WithdrawID,
			"tx_hash":       row.TxHash,
			"log_index":     row.LogIndex,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
			"balance_real":  row.BalanceReal,
			"hex":           row.Hex,
			"create_time":   row.CreateTime,
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

// SQLDeleteTSendEos 删除
func SQLDeleteTSendEos(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_send_eos
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTx 创建
func SQLCreateTTx(ctx context.Context, tx mcommon.DbExeAble, row *DBTTx, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       create_time,
       handle_status,
       handle_msg,
       handle_time,
       org_status,
       org_msg,
       org_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateTTxDuplicate 创建更新
func SQLCreateTTxDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTTx, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       create_time,
       handle_status,
       handle_msg,
       handle_time,
       org_status,
       org_msg,
       org_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateManyTTx 创建多个
func SQLCreateManyTTx(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTx, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTTxDuplicate 创建多个
func SQLCreateManyTTxDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTx, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTxCol 根据id查询
func SQLGetTTxCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTTx, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
WHERE
	id=:id`)

	var row DBTTx
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTTxColKV 根据id查询
func SQLGetTTxColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTTx, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTTx
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTTxCol 根据ids获取
func SQLSelectTTxCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTTx, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTTx
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxColKV 根据ids获取
func SQLSelectTTxColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTTx, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTTx
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTx 更新
func SQLUpdateTTx(ctx context.Context, tx mcommon.DbExeAble, row *DBTTx) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx
SET
    product_id=:product_id,
    tx_id=:tx_id,
    from_address=:from_address,
    to_address=:to_address,
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
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
func SQLDeleteTTx(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTxBtc 创建
func SQLCreateTTxBtc(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtc, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_btc ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
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
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateTTxBtcDuplicate 创建更新
func SQLCreateTTxBtcDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtc, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_btc ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
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
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateManyTTxBtc 创建多个
func SQLCreateManyTTxBtc(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxBtc, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_btc ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTTxBtcDuplicate 创建多个
func SQLCreateManyTTxBtcDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxBtc, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_btc ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTxBtcCol 根据id查询
func SQLGetTTxBtcCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTTxBtc, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc
WHERE
	id=:id`)

	var row DBTTxBtc
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTTxBtcColKV 根据id查询
func SQLGetTTxBtcColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTTxBtc, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTTxBtc
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTTxBtcCol 根据ids获取
func SQLSelectTTxBtcCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTTxBtc, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTTxBtc
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxBtcColKV 根据ids获取
func SQLSelectTTxBtcColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTTxBtc, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTTxBtc
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxBtc 更新
func SQLUpdateTTxBtc(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtc) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
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
func SQLDeleteTTxBtc(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_btc
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTxBtcToken 创建
func SQLCreateTTxBtcToken(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtcToken, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_btc_token ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
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
       handle_at,
       org_status,
       org_msg,
       org_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
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
    :handle_at,
    :org_status,
    :org_msg,
    :org_at
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
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
			"org_status":    row.OrgStatus,
			"org_msg":       row.OrgMsg,
			"org_at":        row.OrgAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTTxBtcTokenDuplicate 创建更新
func SQLCreateTTxBtcTokenDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtcToken, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_btc_token ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
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
       handle_at,
       org_status,
       org_msg,
       org_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
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
    :handle_at,
    :org_status,
    :org_msg,
    :org_at
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
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
			"org_status":    row.OrgStatus,
			"org_msg":       row.OrgMsg,
			"org_at":        row.OrgAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTTxBtcToken 创建多个
func SQLCreateManyTTxBtcToken(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxBtcToken, isIgnore bool) (int64, error) {
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
					row.OrgStatus,
					row.OrgMsg,
					row.OrgAt,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ProductID,
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
					row.OrgStatus,
					row.OrgMsg,
					row.OrgAt,
				},
			)
		}
	}
	var count int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_btc_token ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
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
    handle_at,
    org_status,
    org_msg,
    org_at
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTTxBtcTokenDuplicate 创建多个
func SQLCreateManyTTxBtcTokenDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxBtcToken, updates []string) (int64, error) {
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
					row.OrgStatus,
					row.OrgMsg,
					row.OrgAt,
				},
			)
		}
	} else {
		for _, row := range rows {
			args = append(
				args,
				[]interface{}{
					row.ProductID,
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
					row.OrgStatus,
					row.OrgMsg,
					row.OrgAt,
				},
			)
		}
	}
	var count int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_btc_token ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
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
    handle_at,
    org_status,
    org_msg,
    org_at
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTxBtcTokenCol 根据id查询
func SQLGetTTxBtcTokenCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTTxBtcToken, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_token
WHERE
	id=:id`)

	var row DBTTxBtcToken
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTTxBtcTokenColKV 根据id查询
func SQLGetTTxBtcTokenColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTTxBtcToken, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_token
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTTxBtcToken
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTTxBtcTokenCol 根据ids获取
func SQLSelectTTxBtcTokenCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTTxBtcToken, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTTxBtcToken
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxBtcTokenColKV 根据ids获取
func SQLSelectTTxBtcTokenColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTTxBtcToken, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_token
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTTxBtcToken
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxBtcToken 更新
func SQLUpdateTTxBtcToken(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtcToken) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_btc_token
SET
    product_id=:product_id,
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
    handle_at=:handle_at,
    org_status=:org_status,
    org_msg=:org_msg,
    org_at=:org_at
WHERE
	id=:id`,
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
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
			"org_status":    row.OrgStatus,
			"org_msg":       row.OrgMsg,
			"org_at":        row.OrgAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLDeleteTTxBtcToken 删除
func SQLDeleteTTxBtcToken(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_btc_token
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTxBtcUxto 创建
func SQLCreateTTxBtcUxto(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtcUxto, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_btc_uxto ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
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
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateTTxBtcUxtoDuplicate 创建更新
func SQLCreateTTxBtcUxtoDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtcUxto, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_btc_uxto ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
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
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
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
	return lastID, nil
}

// SQLCreateManyTTxBtcUxto 创建多个
func SQLCreateManyTTxBtcUxto(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxBtcUxto, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_btc_uxto ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTTxBtcUxtoDuplicate 创建多个
func SQLCreateManyTTxBtcUxtoDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxBtcUxto, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_btc_uxto ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
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
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTxBtcUxtoCol 根据id查询
func SQLGetTTxBtcUxtoCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTTxBtcUxto, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_uxto
WHERE
	id=:id`)

	var row DBTTxBtcUxto
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTTxBtcUxtoColKV 根据id查询
func SQLGetTTxBtcUxtoColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTTxBtcUxto, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_uxto
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTTxBtcUxto
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTTxBtcUxtoCol 根据ids获取
func SQLSelectTTxBtcUxtoCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTTxBtcUxto, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTTxBtcUxto
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxBtcUxtoColKV 根据ids获取
func SQLSelectTTxBtcUxtoColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTTxBtcUxto, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_btc_uxto
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTTxBtcUxto
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxBtcUxto 更新
func SQLUpdateTTxBtcUxto(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxBtcUxto) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
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
func SQLDeleteTTxBtcUxto(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_btc_uxto
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTxEos 创建
func SQLCreateTTxEos(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxEos, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_eos ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
       tx_hash,
       log_index,
       from_address,
       to_address,
       memo,
       balance_real,
       create_at,
       handle_status,
       handle_msg,
       handle_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
    :tx_hash,
    :log_index,
    :from_address,
    :to_address,
    :memo,
    :balance_real,
    :create_at,
    :handle_status,
    :handle_msg,
    :handle_at
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"tx_hash":       row.TxHash,
			"log_index":     row.LogIndex,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
			"balance_real":  row.BalanceReal,
			"create_at":     row.CreateAt,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_at":     row.HandleAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateTTxEosDuplicate 创建更新
func SQLCreateTTxEosDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxEos, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_eos ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
       tx_hash,
       log_index,
       from_address,
       to_address,
       memo,
       balance_real,
       create_at,
       handle_status,
       handle_msg,
       handle_at
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
    :tx_hash,
    :log_index,
    :from_address,
    :to_address,
    :memo,
    :balance_real,
    :create_at,
    :handle_status,
    :handle_msg,
    :handle_at
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"tx_hash":       row.TxHash,
			"log_index":     row.LogIndex,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
			"balance_real":  row.BalanceReal,
			"create_at":     row.CreateAt,
			"handle_status": row.HandleStatus,
			"handle_msg":    row.HandleMsg,
			"handle_at":     row.HandleAt,
		},
	)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// SQLCreateManyTTxEos 创建多个
func SQLCreateManyTTxEos(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxEos, isIgnore bool) (int64, error) {
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
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
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
					row.ProductID,
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_eos ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
    tx_hash,
    log_index,
    from_address,
    to_address,
    memo,
    balance_real,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTTxEosDuplicate 创建多个
func SQLCreateManyTTxEosDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxEos, updates []string) (int64, error) {
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
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
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
					row.ProductID,
					row.TxHash,
					row.LogIndex,
					row.FromAddress,
					row.ToAddress,
					row.Memo,
					row.BalanceReal,
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_eos ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
    tx_hash,
    log_index,
    from_address,
    to_address,
    memo,
    balance_real,
    create_at,
    handle_status,
    handle_msg,
    handle_at
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTxEosCol 根据id查询
func SQLGetTTxEosCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTTxEos, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_eos
WHERE
	id=:id`)

	var row DBTTxEos
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTTxEosColKV 根据id查询
func SQLGetTTxEosColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTTxEos, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_eos
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTTxEos
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTTxEosCol 根据ids获取
func SQLSelectTTxEosCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTTxEos, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_eos
WHERE
	id IN (:ids)`)
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTTxEos
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxEosColKV 根据ids获取
func SQLSelectTTxEosColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTTxEos, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_eos
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTTxEos
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxEos 更新
func SQLUpdateTTxEos(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxEos) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_tx_eos
SET
    product_id=:product_id,
    tx_hash=:tx_hash,
    log_index=:log_index,
    from_address=:from_address,
    to_address=:to_address,
    memo=:memo,
    balance_real=:balance_real,
    create_at=:create_at,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_at=:handle_at
WHERE
	id=:id`,
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"tx_hash":       row.TxHash,
			"log_index":     row.LogIndex,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
			"balance_real":  row.BalanceReal,
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

// SQLDeleteTTxEos 删除
func SQLDeleteTTxEos(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_eos
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTTxErc20 创建
func SQLCreateTTxErc20(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxErc20, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_erc20 ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       token_id,
       product_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       create_time,
       handle_status,
       handle_msg,
       handle_time,
       org_status,
       org_msg,
       org_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :token_id,
    :product_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"token_id":      row.TokenID,
			"product_id":    row.ProductID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateTTxErc20Duplicate 创建更新
func SQLCreateTTxErc20Duplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxErc20, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_erc20 ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       token_id,
       product_id,
       tx_id,
       from_address,
       to_address,
       balance_real,
       create_time,
       handle_status,
       handle_msg,
       handle_time,
       org_status,
       org_msg,
       org_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :token_id,
    :product_id,
    :tx_id,
    :from_address,
    :to_address,
    :balance_real,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time,
    :org_status,
    :org_msg,
    :org_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"token_id":      row.TokenID,
			"product_id":    row.ProductID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
	return lastID, nil
}

// SQLCreateManyTTxErc20 创建多个
func SQLCreateManyTTxErc20(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxErc20, isIgnore bool) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_tx_erc20 ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    token_id,
    product_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTTxErc20Duplicate 创建多个
func SQLCreateManyTTxErc20Duplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTTxErc20, updates []string) (int64, error) {
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_tx_erc20 ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    token_id,
    product_id,
    tx_id,
    from_address,
    to_address,
    balance_real,
    create_time,
    handle_status,
    handle_msg,
    handle_time,
    org_status,
    org_msg,
    org_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTTxErc20Col 根据id查询
func SQLGetTTxErc20Col(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTTxErc20, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
WHERE
	id=:id`)

	var row DBTTxErc20
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTTxErc20ColKV 根据id查询
func SQLGetTTxErc20ColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTTxErc20, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTTxErc20
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTTxErc20Col 根据ids获取
func SQLSelectTTxErc20Col(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTTxErc20, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTTxErc20
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTTxErc20ColKV 根据ids获取
func SQLSelectTTxErc20ColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTTxErc20, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_tx_erc20
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTTxErc20
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTTxErc20 更新
func SQLUpdateTTxErc20(ctx context.Context, tx mcommon.DbExeAble, row *DBTTxErc20) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
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
		mcommon.H{
			"id":            row.ID,
			"token_id":      row.TokenID,
			"product_id":    row.ProductID,
			"tx_id":         row.TxID,
			"from_address":  row.FromAddress,
			"to_address":    row.ToAddress,
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
func SQLDeleteTTxErc20(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_tx_erc20
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateTWithdraw 创建
func SQLCreateTWithdraw(ctx context.Context, tx mcommon.DbExeAble, row *DBTWithdraw, isIgnore bool) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_withdraw ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
       out_serial,
       to_address,
       memo,
       symbol,
       balance_real,
       tx_hash,
       create_time,
       handle_status,
       handle_msg,
       handle_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
    :out_serial,
    :to_address,
    :memo,
    :symbol,
    :balance_real,
    :tx_hash,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
)`)
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"out_serial":    row.OutSerial,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
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
	return lastID, nil
}

// SQLCreateTWithdrawDuplicate 创建更新
func SQLCreateTWithdrawDuplicate(ctx context.Context, tx mcommon.DbExeAble, row *DBTWithdraw, updates []string) (int64, error) {
	var lastID int64
	var err error
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_withdraw ( ")
	if row.ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
       product_id,
       out_serial,
       to_address,
       memo,
       symbol,
       balance_real,
       tx_hash,
       create_time,
       handle_status,
       handle_msg,
       handle_time
) VALUES (`)
	if row.ID > 0 {
		query.WriteString("\n:id,")
	}
	query.WriteString(`
    :product_id,
    :out_serial,
    :to_address,
    :memo,
    :symbol,
    :balance_real,
    :tx_hash,
    :create_time,
    :handle_status,
    :handle_msg,
    :handle_time
) `)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	lastID, err = mcommon.DbExecuteLastIDNamedContent(
		ctx,
		tx,
		query.String(),
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"out_serial":    row.OutSerial,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
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
	return lastID, nil
}

// SQLCreateManyTWithdraw 创建多个
func SQLCreateManyTWithdraw(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTWithdraw, isIgnore bool) (int64, error) {
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
					row.Memo,
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
					row.Memo,
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
	query := strings.Builder{}
	query.WriteString("INSERT ")
	if isIgnore {
		query.WriteString("IGNORE ")
	}
	query.WriteString("INTO t_withdraw ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
    out_serial,
    to_address,
    memo,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`)
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLCreateManyTWithdrawDuplicate 创建多个
func SQLCreateManyTWithdrawDuplicate(ctx context.Context, tx mcommon.DbExeAble, rows []*DBTWithdraw, updates []string) (int64, error) {
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
					row.Memo,
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
					row.Memo,
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
	query := strings.Builder{}
	query.WriteString("INSERT INTO t_withdraw ( ")
	if rows[0].ID > 0 {
		query.WriteString("\nid,")
	}
	query.WriteString(`
    product_id,
    out_serial,
    to_address,
    memo,
    symbol,
    balance_real,
    tx_hash,
    create_time,
    handle_status,
    handle_msg,
    handle_time
) VALUES
    %s`)
	updatesLen := len(updates)
	lastUpdateIndex := updatesLen - 1
	if updatesLen > 0 {
		query.WriteString("ON DUPLICATE KEY UPDATE\n")
		for i, update := range updates {
			query.WriteString(update)
			query.WriteString("=VALUES(")
			query.WriteString(update)
			query.WriteString(")")
			if i != lastUpdateIndex {
				query.WriteString(",\n")
			} else {
				query.WriteString("\n")
			}
		}
	}
	count, err = mcommon.DbExecuteCountManyContent(
		ctx,
		tx,
		query.String(),
		len(rows),
		args...,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// SQLGetTWithdrawCol 根据id查询
func SQLGetTWithdrawCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, id int64) (*DBTWithdraw, error) {
	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
WHERE
	id=:id`)

	var row DBTWithdraw
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		mcommon.H{
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

// SQLGetTWithdrawColKV 根据id查询
func SQLGetTWithdrawColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}) (*DBTWithdraw, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	var row DBTWithdraw
	ok, err := mcommon.DbGetNamedContent(
		ctx,
		tx,
		&row,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return &row, nil
}

// SQLSelectTWithdrawCol 根据ids获取
func SQLSelectTWithdrawCol(ctx context.Context, tx mcommon.DbExeAble, cols []string, ids []int64, orderBys []string, limits []int64) ([]*DBTWithdraw, error) {
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
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}
	var rows []*DBTWithdraw
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		mcommon.H{
			"ids": ids,
		},
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLSelectTWithdrawColKV 根据ids获取
func SQLSelectTWithdrawColKV(ctx context.Context, tx mcommon.DbExeAble, cols []string, keys []string, values []interface{}, orderBys []string, limits []int64) ([]*DBTWithdraw, error) {
	keysLen := len(keys)
	if keysLen != len(values) {
		return nil, fmt.Errorf("value len error")
	}

	query := strings.Builder{}
	query.WriteString("SELECT\n")
	query.WriteString(strings.Join(cols, ",\n"))
	query.WriteString(`
FROM
	t_withdraw
`)
	if len(keys) > 0 {
		query.WriteString("WHERE\n")
	}
	argMap := mcommon.H{}
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return nil, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(" )")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}
	if len(orderBys) > 0 {
		query.WriteString("\nORDER BY\n")
		query.WriteString(strings.Join(orderBys, ",\n"))
		query.WriteString("\n")
	}
	if len(limits) == 1 {
		query.WriteString(fmt.Sprintf("LIMIT %d", limits[0]))
	}
	if len(limits) == 2 {
		query.WriteString(fmt.Sprintf("LIMIT %d,%d", limits[0], limits[1]))
	}

	var rows []*DBTWithdraw
	err := mcommon.DbSelectNamedContent(
		ctx,
		tx,
		&rows,
		query.String(),
		argMap,
	)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SQLUpdateTWithdraw 更新
func SQLUpdateTWithdraw(ctx context.Context, tx mcommon.DbExeAble, row *DBTWithdraw) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`UPDATE
	t_withdraw
SET
    product_id=:product_id,
    out_serial=:out_serial,
    to_address=:to_address,
    memo=:memo,
    symbol=:symbol,
    balance_real=:balance_real,
    tx_hash=:tx_hash,
    create_time=:create_time,
    handle_status=:handle_status,
    handle_msg=:handle_msg,
    handle_time=:handle_time
WHERE
	id=:id`,
		mcommon.H{
			"id":            row.ID,
			"product_id":    row.ProductID,
			"out_serial":    row.OutSerial,
			"to_address":    row.ToAddress,
			"memo":          row.Memo,
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
func SQLDeleteTWithdraw(ctx context.Context, tx mcommon.DbExeAble, id int64) (int64, error) {
	count, err := mcommon.DbExecuteCountNamedContent(
		ctx,
		tx,
		`DELETE
FROM
	t_withdraw
WHERE
	id=:id`,
		mcommon.H{
			"id": id,
		},
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}
