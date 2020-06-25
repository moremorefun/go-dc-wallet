package main

import (
	"bytes"
	"context"
	"go-dc-wallet/app"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/schemalex/schemalex/diff"
)

func main() {
	app.EnvCreate()
	defer app.EnvDestroy()

	var dbSQLs []string
	for _, tableName := range model.TableNames {
		var row struct {
			TableName string `db:"Table"`
			TableSQL  string `db:"Create Table"`
		}
		ok, err := hcommon.DbGetNamedContent(
			context.Background(),
			app.DbCon,
			&row,
			`SHOW CREATE TABLE `+tableName,
			gin.H{},
		)
		if err != nil {
			if strings.Contains(err.Error(), "doesn't exist") {
				continue
			}
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			return
		}
		if ok {
			dbSQLs = append(dbSQLs, row.TableSQL+";")
		}
	}
	// 原始sql
	dbSQL := strings.Join(dbSQLs, "\n")
	// 目的sql
	toSQL, err := ioutil.ReadFile("init/dc-wallet.sql")
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	sqlDiff := new(bytes.Buffer)
	err = diff.Strings(sqlDiff, dbSQL, string(toSQL), diff.WithTransaction(true))
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	r, _ := regexp.Compile(`AUTO_INCREMENT\s*=\s*(\d)*\s*,`)
	sqlDiffWithoutInc := r.ReplaceAllStringFunc(sqlDiff.String(), func(s string) string {
		return ""
	})
	hcommon.Log.Debugf("sql diff: %s", sqlDiffWithoutInc)
}
