package myrunner

import "github.com/leeprince/gopublic/errors"

//任务执行超时
var ErrTimeOut = errors.New("run time out")

//任务执行中断
var ErrInterruput = errors.New("run interruput")
