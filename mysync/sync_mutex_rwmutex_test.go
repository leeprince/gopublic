package mysync

import (
    "fmt"
    "sync"
    "testing"
)

/**
 * @Author: prince.lee <leeprince@foxmail.com>
 * @Date:   2022/1/3 下午3:43
 * @Desc:   当程序中就一个线程的时候，是不需要加锁的，但是通常实际的代码不会只是单线程，所以这个时候就需要用到锁了，那么关于锁的使用场景主要涉及到哪些呢？
                 1. 多个线程在读相同的数据时
                 2. 多个线程在写相同的数据时
                 3. 同一个资源，有读又有写
 
 */

// 互斥锁是一种常用的控制共享资源访问的方法，它能够保证同时只有一个 goroutine 可以访问到共享资源（同一个时刻只有一个线程能够拿到锁）
func TestMutex(t *testing.T) {
    var (
        count int
        lock  sync.Mutex // 修改代码，在累加的地方添加互斥锁，就能保证我们每次得到的结果都是想要的值
    )
    
    wg := sync.WaitGroup{}
    wg.Add(2)
    for i := 0; i < 2; i++ {
        go func() {
            defer wg.Done()
            for i := 100000; i > 0; i-- {
                lock.Lock()
                count++
                lock.Unlock()
            }
            fmt.Println(count)
        }()
    }
    wg.Wait()
    fmt.Printf("最后的结果：%d \n", count) // 等待子线程全部结束
}

// 在读多写少的环境中，可以优先使用读写互斥锁（sync.RWMutex），它比互斥锁更加高效。sync 包中的 RWMutex 提供了读写互斥锁的封装。
// 读写锁分为：读锁和写锁
//    如果设置了一个写锁，那么其它读的线程以及写的线程都拿不到锁，这个时候，与互斥锁的功能相同
//    如果设置了一个读锁，那么其它写的线程是拿不到锁的，但是其它读的线程是可以拿到锁
func TestRWMutex(t *testing.T) {
    var (
        count int
        lock  sync.RWMutex // 修改代码，在累加的地方添加互斥锁，就能保证我们每次得到的结果都是想要的值
    )
    
    wg := sync.WaitGroup{}
    wg.Add(2)
    for i := 0; i < 2; i++ {
        go func() {
            defer wg.Done()
            for i := 100000; i > 0; i-- {
                lock.RLock()
                count++
                lock.RUnlock()
                
                // 读锁，非预期
                // lock.RLock()
                // count++
                // lock.RUnlock()
            }
            fmt.Println(count)
        }()
    }
    wg.Wait()
    fmt.Printf("最后的结果：%d \n", count) // 等待子线程全部结束
}
