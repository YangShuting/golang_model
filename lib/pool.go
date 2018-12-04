package lib

import (
	"errors"
	"io"
	"log"
	"sync"
)

//实现一个有缓冲通道的资源池，可以管理在任意多个 goroutine之间的资源共享，比如网络连接和数据库连接等。
//每个 goroutine 可以向资源池里申请资源，然后使用完之后放回资源池里。
type Pool struct {
	m       sync.Mutex                //互斥锁，这主要是用来保证在多个goroutine访问资源时，池内的值是安全的。
	res     chan io.Closer            //有缓冲的通道，用来保存共享的资源
	factory func() (io.Closer, error) //当需要一个新的资源时，可以通过这个函数创建
	closed  bool                      //资源池是否被关闭，如果被关闭的话，再访问是会有错误的
}

var ErrorPoolClosed = errors.New("资源已经关闭")

//创建一个资源池
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("资源池太小了")
	}
	return &Pool{
		res:     make(chan io.Closer, size),
		factory: fn,
	}, nil
}

//从资源池里获取一个资源
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.res:
		log.Println("Acquire：共享资源")
		if !ok {
			return nil, ErrorPoolClosed
		}
		return r, nil
	default:
		log.Println("Aquire: 新生资源")
		return p.factory()
	}
}

//关闭资源池
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	if p.closed {
		return
	}
	p.closed = true
	close(p.res)
	for r := range p.res {
		r.Close()
	}
}

//然后释放资源池里的资源
func (p *Pool) Release(r io.Closer) {
	//保证操作是安全的
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		//如果 Close()操作同时在进行，那么能保证只有其中一个在操作。
		r.Close()
		return
	}

	select {
	case p.res <- r:
		log.Println("资源池释放到池子里了")
	default:
		log.Println("资源满了，释放这个资源吧")
		r.Close()
	}
}
