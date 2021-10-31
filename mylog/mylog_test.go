package mylog

import (
	"fmt"
	"testing"
	"time"

	"github.com/leeprince/gopublic/dev"

	"github.com/leeprince/gopublic/errors"
)

func TestStdLog(t *testing.T) {
	Info("aaaa")
	Debug("bbbb")
	Error(fmt.Errorf("nottttt"))
	fmt.Println(TraceError(errors.New("wwww")))
}

type aaa struct {
	Act   string
	Begin int64
}
type ttt struct {
	Act   string
	Begin int64
	Cat   *aaa
}

func TestZapLog(t *testing.T) {
	dev.SetService("xxjwxc")
	//	dev.OnSetDev(true)
	SetLog(GetDefaultZap())
	// log.Printf("%#v", &ttt{
	// 	Act:   "====001===",
	// 	Begin: time.Now().Unix(),
	// 	Cat: &aaa{
	// 		Act:   "----002----",
	// 		Begin: time.Now().Unix(),
	// 	},
	// })
	// return

	Info(&ttt{
		Act:   "====001===",
		Begin: time.Now().Unix(),
		Cat: &aaa{
			Act:   "----002----",
			Begin: time.Now().Unix(),
		},
	})
	Info("aaaa")
	Debug("bbbb")
	Error(fmt.Errorf("nottttt"))
	fmt.Println(TraceError(errors.New("wwww")))

}
