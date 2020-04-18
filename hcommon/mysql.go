package hcommon

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"time"

	// 导入mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DbExeAble 数据操作接口
type DbExeAble interface {
	Rebind(string) string
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

var isShowSQL bool

// DbCreate 创建数据库
func DbCreate(dataSourceName string, showSQL bool) *sqlx.DB {
	isShowSQL = showSQL

	var err error
	var db *sqlx.DB

	db, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		Log.Fatalf("db connect error: %s", err.Error())
		return nil
	}

	count := runtime.NumCPU()*20 + 1
	db.SetMaxOpenConns(count)
	db.SetMaxIdleConns(count)
	db.SetConnMaxLifetime(1 * time.Hour)

	err = db.Ping()
	if err != nil {
		Log.Fatalf("db ping error: %s", err.Error())
		return nil
	}
	return db
}

// DbExecuteCountManyContent 返回执行个数
func DbExecuteCountManyContent(ctx context.Context, tx DbExeAble, query string, n int, args ...interface{}) (int64, error) {
	var err error
	insertArgs := strings.Repeat("(?),", n)
	insertArgs = strings.TrimSuffix(insertArgs, ",")
	query = fmt.Sprintf(query, insertArgs)
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
	}
	ret, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return 0, err
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DbExecuteLastIDNamedContent 执行并返回lastID
func DbExecuteLastIDNamedContent(ctx context.Context, tx DbExeAble, query string, argMap map[string]interface{}) (int64, error) {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
	}
	ret, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return 0, err
	}
	lastID, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// DbExecuteCountNamedContent 返回执行个数
func DbExecuteCountNamedContent(ctx context.Context, tx DbExeAble, query string, argMap map[string]interface{}) (int64, error) {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
	}
	ret, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return 0, err
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DbGetNamedContent 返回单个元素
func DbGetNamedContent(ctx context.Context, tx DbExeAble, dest interface{}, query string, argMap map[string]interface{}) (bool, error) {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return false, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return false, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
	}
	err = tx.GetContext(
		ctx,
		dest,
		query,
		args...,
	)
	if err == sql.ErrNoRows {
		// 没有元素
		return false, nil
	}
	if err != nil {
		// 执行错误
		return false, err
	}
	return true, nil
}

// DbSelectNamedContent 返回列表
func DbSelectNamedContent(ctx context.Context, tx DbExeAble, dest interface{}, query string, argMap map[string]interface{}) error {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
	}
	err = tx.SelectContext(
		ctx,
		dest,
		query,
		args...,
	)
	if err == sql.ErrNoRows {
		// 没有元素
		return nil
	}
	if err != nil {
		// 执行错误
		return err
	}
	return nil
}
