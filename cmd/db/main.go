package main

import (
	"bytes"
	"context"
	"fmt"
	"go-dc-wallet/hcommon"
	"go-dc-wallet/model"
	"go-dc-wallet/xenv"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/schemalex/schemalex/diff"
)

func main() {
	xenv.EnvCreate()
	defer xenv.EnvDestroy()

	var dbSQLs []string
	for _, tableName := range model.TableNames {
		var row struct {
			TableName string `db:"Table"`
			TableSQL  string `db:"Create Table"`
		}
		ok, err := hcommon.DbGetNamedContent(
			context.Background(),
			xenv.DbCon,
			&row,
			`SHOW CREATE TABLE `+tableName,
			gin.H{},
		)
		if err != nil {
			if strings.Contains(err.Error(), "doesn't exist") {
				continue
			}
			hcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
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
		hcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	sqlDiff := new(bytes.Buffer)
	err = diff.Strings(sqlDiff, dbSQL, string(toSQL), diff.WithTransaction(true))
	if err != nil {
		hcommon.Log.Fatalf("err: [%T] %s", err, err.Error())
	}
	// 替换 AUTO_INCREMENT
	r, _ := regexp.Compile(`AUTO_INCREMENT\s*=\s*(\d)*\s*,`)
	sqlDiffWithoutInc := r.ReplaceAllStringFunc(sqlDiff.String(), func(s string) string {
		return ""
	})
	fmt.Printf("%s\n", sqlDiffWithoutInc)
}
