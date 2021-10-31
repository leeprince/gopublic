package timerDeal

import (
	"log"

	"github.com/leeprince/gopublic/tools"

	"time"

	"github.com/leeprince/gopublic/mydef"
	"github.com/leeprince/gopublic/mylog"
)

/*
	主要为定时任务
*/

/*
超时回调
t:超时时间
fun:回调接口
args:回调接口传入的参数
*/
func OnDealTimeOut(t time.Duration, fun mydef.ParamsCallFunc, parms ...interface{}) {
	go func() {
		ticker := time.NewTicker(t)
		<-ticker.C
		mylog.Debug("timer 执行.....")
		//以下为定时执行的操作
		fun(parms...)
	}()
}

/*
	每月事件
	day : 几号
	hour, min, sec : 几点(当天的0点偏移秒数)
	callback : 时间回调
*/
func OnPeMonth(day int, hour, min, sec int, callback func()) {
	go func() {
		for {
			next := time.Now().AddDate(0, 1, 0)
			next = time.Date(next.Year(), next.Month(), day, hour, min, sec, 0, next.Location())
			t := time.NewTimer(next.Sub(time.Now()))
			log.Println("next time callback:", next)
			<-t.C
			callback()
		}
	}()
}

/*
	每天事件
	hour, min, sec : 几点(当天的0点偏移秒数)
	callback : 时间回调
*/
func OnPeDay(hour, min, sec int, callback func()) {
	go func() {
		next := time.Now()
		for {
			next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, sec, 0, next.Location())
			mylog.Infof("next pe day on:%v", tools.GetTimeStr(next))
			t := time.NewTimer(next.Sub(time.Now()))
			log.Println("next time callback:", next)
			<-t.C
			callback()
			next = time.Now().AddDate(0, 0, 1) // 下一天
		}
	}()
}
