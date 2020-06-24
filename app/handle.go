package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go-dc-wallet/app/model"
	"go-dc-wallet/hcommon"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/parnurzeal/gorequest"
)

// CheckDoNotify 检测发送回调
func CheckDoNotify() {
	lockKey := "CheckDoNotify"
	ok, err := GetLock(
		context.Background(),
		DbCon,
		lockKey,
	)
	if err != nil {
		hcommon.Log.Warnf("GetLock err: [%T] %s", err, err.Error())
		return
	}
	if !ok {
		return
	}
	defer func() {
		err := ReleaseLock(
			context.Background(),
			DbCon,
			lockKey,
		)
		if err != nil {
			hcommon.Log.Warnf("ReleaseLock err: [%T] %s", err, err.Error())
			return
		}
	}()

	// 初始化的
	initNotifyRows, err := SQLSelectTProductNotifyColByStatusAndTime(
		context.Background(),
		DbCon,
		[]string{
			model.DBColTProductNotifyID,
			model.DBColTProductNotifyURL,
			model.DBColTProductNotifyMsg,
		},
		NotifyStatusInit,
		time.Now().Unix(),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	// 错误的
	delayNotifyRows, err := SQLSelectTProductNotifyColByStatusAndTime(
		context.Background(),
		DbCon,
		[]string{
			model.DBColTProductNotifyID,
			model.DBColTProductNotifyURL,
			model.DBColTProductNotifyMsg,
		},
		NotifyStatusFail,
		time.Now().Add(-time.Minute*10).Unix(),
	)
	if err != nil {
		hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
		return
	}
	initNotifyRows = append(initNotifyRows, delayNotifyRows...)

	for _, initNotifyRow := range initNotifyRows {
		gresp, body, errs := gorequest.New().Post(initNotifyRow.URL).Timeout(time.Second * 30).Send(initNotifyRow.Msg).End()
		if errs != nil {
			hcommon.Log.Errorf("err: [%T] %s", errs[0], errs[0].Error())
			_, err = SQLUpdateTProductNotifyStatusByID(
				context.Background(),
				DbCon,
				&model.DBTProductNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: NotifyStatusFail,
					HandleMsg:    errs[0].Error(),
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
		if gresp.StatusCode != http.StatusOK {
			// 状态错误
			hcommon.Log.Errorf("req status error: %d", gresp.StatusCode)
			_, err = SQLUpdateTProductNotifyStatusByID(
				context.Background(),
				DbCon,
				&model.DBTProductNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: NotifyStatusFail,
					HandleMsg:    fmt.Sprintf("http status: %d", gresp.StatusCode),
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
		resp := gin.H{}
		err = json.Unmarshal([]byte(body), &resp)
		if err != nil {
			hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			_, err = SQLUpdateTProductNotifyStatusByID(
				context.Background(),
				DbCon,
				&model.DBTProductNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: NotifyStatusFail,
					HandleMsg:    body,
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
		_, ok := resp["error"]
		if ok {
			// 处理成功
			_, err = SQLUpdateTProductNotifyStatusByID(
				context.Background(),
				DbCon,
				&model.DBTProductNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: NotifyStatusPass,
					HandleMsg:    body,
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
		} else {
			//hcommon.Log.Errorf("no error in resp")
			_, err = SQLUpdateTProductNotifyStatusByID(
				context.Background(),
				DbCon,
				&model.DBTProductNotify{
					ID:           initNotifyRow.ID,
					HandleStatus: NotifyStatusFail,
					HandleMsg:    body,
					UpdateTime:   time.Now().Unix(),
				},
			)
			if err != nil {
				hcommon.Log.Errorf("err: [%T] %s", err, err.Error())
			}
			continue
		}
	}
}
