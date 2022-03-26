package mysync

import (
    "errors"
    "fmt"
    "golang.org/x/sync/errgroup"
    "math/rand"
    "sync"
    "testing"
    "time"
)

/**
 * @Author: prince.lee <leeprince@foxmail.com>
 * @Date:   2022/1/3 下午3:11
 * @Desc:   https://mp.weixin.qq.com/s?__biz=MzAwNjMxMTgwNw==&mid=2247490256&idx=1&sn=4993c8ce4f80045dc11162f0f801d128&chksm=9b0e00c0ac7989d65c3676f8403525b83b0c5b74e71b13c5d6a5ad2bbdcf6c8b5ec829a0ac04&scene=90&subscene=93&sessionid=1640924932&clicktime=1640924938&enterid=1640924938&ascene=56&devicetype=android-29&version=2800105b&nettype=cmnet&abtest_cookie=AAACAA%3D%3D&lang=zh_CN&exportkey=AcpV08PwTTNljDmLw99dTrs%3D&pass_ticket=%2BKCCfEhtv3mUHDPs00hJgbX30HnIC7FPNwKaYdDfBL7AYqmqJE0F7E7prT9UJvw0&wx_header=1
 *          sync.Cond：为等待 / 通知场景下的并发问题提供支持。
 *          Cond 通常应用于等待某个条件的一组 goroutine，等条件变为 true 的时候，其中一个 goroutine 或者所有的 goroutine 都会被唤醒执行。
 *               func NewCond(l Locker) *Cond {} // 创建一个 cond
 *               func (c *Cond) Wait() {}        // 阻塞，等待唤醒
 *               func (c *Cond) Signal() {}      // 唤醒一个等待者
 *               func (c *Cond) Broadcast() {}   // 唤醒所有等待者
 *          场景：场景是百米赛跑，10个运动员，进场以后做热身运动，运动员热身完成后示意裁判，10个运动员都热身完成，裁判发令起跑
 */

// 这里你可能会说，使用 sync.WaitGroup{} 或 channel 也可以实现，甚至比 cond 的实现还要简单，的确如此，这也从侧面说明 cond 的应用场景少之又少。
// sync.WaitGroup{} 或 channel 这种并发原语适用的情况时，等待者只有一个，如果等待者有多个，cond 比较擅长。
// 一个裁判: sync.Cond 实现
func TestCond(t *testing.T) {
    c := sync.NewCond(&sync.Mutex{})
    var readyCnt int
    athletesNum := 100 // 运动员数
    
    for i := 0; i < athletesNum; i++ {
        go func(i int) {
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", i)
            time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
            
            // 热身结束，加锁更改等待条件
            c.L.Lock()
            readyCnt++
            fmt.Printf("运动员#%d 已热身结束, 运动员总数：>>>>>>: %d  \n", i, readyCnt)
            c.L.Unlock()
            
            c.Signal() // 示意裁判员
        }(i)
    }
    
    // 注意：c.wait() 里面有：c.L.Lock() 的步骤，所以有可能未能及时获取到锁，可能被readyCnt++前c.L.Lock()获取锁成功，所以打印出来的 readyCnt 可能是不连续的
    var realNum int
    c.L.Lock()
    for readyCnt != athletesNum { // 每次 c.Signal() 都会唤醒一次，唤醒 10 次才能开始比赛
        realNum++
        fmt.Println("======准备进入堵塞........readyCnt:", readyCnt, "-实际获取锁唤醒次数:", realNum)
        c.Wait() // c.Wait() 调用后，会阻塞在这里，直到被唤醒。调用 Wait() 时，它会把当前 goroutine 放入等待队列，然后解锁，将自己阻塞等待唤醒，当有其它 goroutine 执行了唤醒操作时，会先获取锁，然后执行 Wait 后面的代码。
        fmt.Printf("======有运动员热身结束，裁判员被唤醒一次.readyCnt: %d \n", readyCnt)
    }
    c.L.Unlock()
    
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
}

// 一个裁判: sync.Cond + sync.Mutex 实现。c.L.Lock() => mut.Lock()
func TestCond1(t *testing.T) {
    c := sync.NewCond(&sync.Mutex{})
    mut := sync.Mutex{}
    var readyCnt int
    athletesNum := 100 // 运动员数
    
    for i := 0; i < athletesNum; i++ {
        go func(i int) {
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", i)
            time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
            
            // 热身结束，加锁更改等待条件
            // c.L.Lock()
            mut.Lock()
            readyCnt++
            fmt.Printf("运动员#%d 已热身结束, 运动员总数：>>>>>>: %d  \n", i, readyCnt)
            mut.Unlock()
            // c.L.Unlock()
            
            c.Signal() // 示意裁判员
        }(i)
    }
    
    // 注意：c.wait() 里面有：c.L.Lock() 的步骤，所以有可能未能及时获取到锁，可能被readyCnt++前c.L.Lock()获取锁成功，所以打印出来的 readyCnt 可能是不连续的
    var realNum int
    c.L.Lock()
    for readyCnt != athletesNum { // 每次 c.Signal() 都会唤醒一次，唤醒 10 次才能开始比赛
        realNum++
        fmt.Println("======准备进入堵塞........readyCnt:", readyCnt, "-实际获取锁唤醒次数:", realNum)
        c.Wait() // c.Wait() 调用后，会阻塞在这里，直到被唤醒。调用 Wait() 时，它会把当前 goroutine 放入等待队列，然后解锁，将自己阻塞等待唤醒，当有其它 goroutine 执行了唤醒操作时，会先获取锁，然后执行 Wait 后面的代码。
        fmt.Printf("======有运动员热身结束，裁判员被唤醒一次.readyCnt: %d \n", readyCnt)
    }
    c.L.Unlock()
    
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
}

// 一个裁判: sync.Cond 实现。c.L.Unlock() 放到  c.Signal() 之后
func TestCond2(t *testing.T) {
    c := sync.NewCond(&sync.Mutex{})
    var readyCnt int
    athletesNum := 100 // 运动员数
    
    for i := 0; i < athletesNum; i++ {
        go func(i int) {
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", i)
            time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
            
            // 热身结束，加锁更改等待条件
            c.L.Lock()
            
            readyCnt++
            fmt.Printf("运动员#%d 已热身结束, 运动员总数：>>>>>>: %d  \n", i, readyCnt)
            
            c.Signal() // 示意裁判员
            fmt.Printf("运动员#%d 示意裁判完成 \n", i)
            
            c.L.Unlock()
        }(i)
    }
    
    // 注意：c.wait() 里面有：c.L.Lock() 的步骤，所以有可能未能及时获取到锁，可能被readyCnt++前c.L.Lock()获取锁成功，所以打印出来的 readyCnt 可能是不连续的
    var realNum int
    c.L.Lock()
    for readyCnt != athletesNum { // 每次 c.Signal() 都会唤醒一次，唤醒 10 次才能开始比赛
        realNum++
        fmt.Println("======准备进入堵塞........readyCnt:", readyCnt, "-实际获取锁唤醒次数:", realNum)
        c.Wait() // c.Wait() 调用后，会阻塞在这里，直到被唤醒。调用 Wait() 时，它会把当前 goroutine 放入等待队列，然后解锁，将自己阻塞等待唤醒，当有其它 goroutine 执行了唤醒操作时，会先获取锁，然后执行 Wait 后面的代码。
        fmt.Printf("======有运动员热身结束，裁判员被唤醒一次.readyCnt: %d \n", readyCnt)
    }
    c.L.Unlock()
    
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
}

// 一个裁判：sync.WaitGroup 实现
func TestWaitGroup(t *testing.T) {
    var readyCnt int
    athletesNum := 100 // 运动员数
    
    mut := sync.Mutex{}
    
    wg := sync.WaitGroup{}
    wg.Add(athletesNum)
    for ii := 0; ii < athletesNum; ii++ {
        go func(ii int) {
            defer wg.Done()
            
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", ii)
            time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
            
            // 热身结束，加锁更改等待条件
            mut.Lock()
            readyCnt++
            fmt.Println("已热身结束的运动员总数：>>>>>>", readyCnt)
            mut.Unlock()
            
            fmt.Printf("运动员#%d 热身结束\n", ii)
        }(ii)
    }
    wg.Wait()
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
}

// 一个裁判："golang.org/x/sync/errgroup" 包 实现
func TestErrGroup(t *testing.T) {
    var readyCnt int
    athletesNum := 100 // 运动员数
    
    mut := sync.Mutex{}
    
    errg := errgroup.Group{}
    for ii := 0; ii < athletesNum; ii++ {
        ii := ii
        errg.Go(func() error {
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", ii)
            time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
            
            // 热身结束，加锁更改等待条件
            mut.Lock()
            readyCnt++
            fmt.Println("已热身结束的运动员总数：>>>>>>", readyCnt)
            mut.Unlock()
            
            fmt.Printf("运动员#%d 热身结束\n", ii)
            
            cTime := time.Now().UnixNano()
            errMsg := fmt.Sprintf("运动员#%d 手动返回错误：#%d time:%d \n", ii, readyCnt, cTime)
            
            if readyCnt <= 4 {
                fmt.Println("++", errMsg)
                return errors.New(errMsg)
            }
            
            return nil
        })
    }
    err := errg.Wait()
    if err != nil {
        fmt.Println("--------------------------err:", err)
    }
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
}

// 一个裁判：Channel 实现
func TestChannel(t *testing.T) {
    var readyCnt int
    athletesNum := 10 // 运动员数
    mut := sync.Mutex{}
    
    ready := make(chan bool)
    // ready := make(chan bool, 10)
    done := make(chan bool)
    
    for i := 0; i < athletesNum; i++ {
        go func(i int) {
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", i)
            time.Sleep(time.Duration(rand.Int63n(1)) * time.Second)
    
            // 热身结束，加锁更改等待条件
            // 解锁的动作要放到通道的后面，否则会出现，readyCnt执行2次或以上自增，导致1.「<- ready:」打印出来的readyCnt会有重复。2. readyCnt == 10并且「if readyCnt == athletesNum {」并发执行，导致执行多次「done<- true」
            mut.Lock()
            readyCnt++
            // fmt.Printf("运动员#%d 已热身结束, 运动员总数：>>>>>>: %d  \n", i, readyCnt)
            
            ready <- true
            mut.Unlock()
    
            // fmt.Printf("运动员#%d 示意裁判完成 \n", i)
    
            if readyCnt == athletesNum {
                done<- true
            }
            
        }(i)
    }
    
    for {
        select {
        case <- ready:
            fmt.Printf("======有运动员热身结束，裁判员被唤醒一次.readyCnt: %d \n", readyCnt)
        case <- done:
            fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
        default:
    
        }
    }
}

// 两个裁判
func TestCond20(t *testing.T) {
    c := sync.NewCond(&sync.Mutex{})
    var readyCnt int
    athletesNum := 10 // 运动员数
    
    for i := 0; i < athletesNum; i++ {
        go func(i int) {
            // 模拟热身
            fmt.Printf("运动员#%d 热身中...\n", i)
            time.Sleep(time.Duration(rand.Int63n(2)) * time.Second)
            
            // 热身结束，加锁更改等待条件
            fmt.Printf("运动员#%d 热身结束，获取锁之前\n", i)
            c.L.Lock()
            readyCnt++
            fmt.Println("已热身结束的运动员总数：>>>>>>", readyCnt)
            c.L.Unlock()
            
            fmt.Printf("运动员#%d 热身结束，已解锁\n", i)
            c.Broadcast() // 示意所有裁判员
        }(i)
    }
    
    wg := sync.WaitGroup{}
    cpNum := 2 // 裁判数
    wg.Add(cpNum)
    for cp := 0; cp < cpNum; cp++ {
        go func(cp int) {
            defer wg.Done()
            fmt.Printf("======裁判员##%d 获取锁之前，准备进入堵塞........readyCnt: %d \n", cp, readyCnt)
            // 注意：c.wait() 里面有：c.L.Lock() 的步骤，所以有可能未能及时获取到锁，可能被readyCnt++前c.L.Lock()获取锁成功，所以打印出来的 readyCnt 可能是不连续的
            c.L.Lock()
            for readyCnt != athletesNum { // 每次 c.Signal() 都会唤醒一次，唤醒 10 次才能开始比赛
                fmt.Printf("======裁判员##%d 准备进入堵塞........readyCnt: %d \n", cp, readyCnt)
                c.Wait() // c.Wait() 调用后，会阻塞在这里，直到被唤醒。调用 Wait() 时，它会把当前 goroutine 放入等待队列，然后解锁，将自己阻塞等待唤醒，当有其它 goroutine 执行了唤醒操作时，会先获取锁，然后执行 Wait 后面的代码。
                fmt.Printf("======有运动员热身结束，裁判员 #%d 被唤醒一次.readyCnt: %d \n", cp, readyCnt)
            }
            c.L.Unlock()
        }(cp)
    }
    wg.Wait()
    
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>所有运动员都热身结束。比赛开始，3，2，1, ......")
}
