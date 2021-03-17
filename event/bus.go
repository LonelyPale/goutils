package event

import (
	"sync"

	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultLogger     Logger = log.New()
	ProcessorPoolSize        = 100 //协程池和订阅管道缓冲区的大小
)

type Logger interface {
	Error(args ...interface{})
}

// Bus 存储有关订阅者感兴趣的特定主题的信息
type Bus struct {
	subscribers map[string]Chans
	mu          sync.RWMutex //mu锁subscribers
	pools       []*ants.PoolWithFunc
	lock        sync.Mutex //lock锁pools
}

func NewBus() *Bus {
	return &Bus{
		subscribers: map[string]Chans{},
		pools:       []*ants.PoolWithFunc{},
	}
}

// 发布
func (eb *Bus) Publish(typ string, data interface{}) {
	eb.mu.RLock()
	if chans, found := eb.subscribers[typ]; found {
		// 这样做是因为切片引用相同的数组，即使它们是按值传递的
		// 因此我们正在使用我们的元素创建一个新切片，从而正确地保持锁定
		channels := append(Chans{}, chans...)
		go func(event Event, eventChans Chans) {
			for _, ch := range eventChans {
				ch <- event
			}
		}(Event{Type: typ, Data: data}, channels)
	}
	eb.mu.RUnlock()
}

// 订阅
func (eb *Bus) Subscribe(typ string, ch Chan) {
	eb.mu.Lock()
	if chans, found := eb.subscribers[typ]; found {
		eb.subscribers[typ] = append(chans, ch)
	} else {
		eb.subscribers[typ] = append([]Chan{}, ch)
	}
	eb.mu.Unlock()
}

// 订阅方法
func (eb *Bus) SubscribeFunc(typ string, fun HandlerFunc) error {
	pool, err := ants.NewPoolWithFunc(ProcessorPoolSize, func(i interface{}) {
		fun(i.(Event))
	})
	if err != nil {
		return err
	}

	eb.lock.Lock()
	eb.pools = append(eb.pools, pool)
	eb.lock.Unlock()

	eventChan := make(Chan, ProcessorPoolSize)
	eb.Subscribe(typ, eventChan)

	//接收事件
	go func() {
		for event := range eventChan {
			if err := pool.Invoke(event); err != nil {
				DefaultLogger.Error(err)
			}
		}
	}()

	return nil
}

// 安全释放协程池
func (eb *Bus) Close() {
	eb.lock.Lock()
	defer eb.lock.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(eb.pools))
	for _, pool := range eb.pools {
		go func() {
			defer wg.Done()
			pool.Release()
		}()
	}
	wg.Wait()
}
