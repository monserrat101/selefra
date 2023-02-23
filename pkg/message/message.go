package message

import (
	"github.com/selefra/selefra-utils/pkg/reflect_util"
	"sync"
)

// Channel 用于链接多个Channel，协调树形调用关系时的消息传递
type Channel[Message any] struct {

	// 当前的channel
	channel chan Message

	// 子channel的控制
	subChannelWg *sync.WaitGroup
	selfWg       *sync.WaitGroup

	// 关闭时的回调
	closeCallbackFunc func()

	// 当前信道处理消息
	consumerFunc func(index int, message Message)
}

func NewChannel[Message any](consumerFunc func(index int, message Message), buffSize ...int) *Channel[Message] {

	// can have buff
	var channel chan Message
	if len(buffSize) != 0 {
		channel = make(chan Message, buffSize[0])
	} else {
		channel = make(chan Message)
	}

	x := &Channel[Message]{
		channel:      channel,
		subChannelWg: &sync.WaitGroup{},
		selfWg:       &sync.WaitGroup{},
		consumerFunc: consumerFunc,
	}

	x.selfWg.Add(1)
	go func() {

		// channel消费者退出时说明被关闭了，则触发关闭时的回调事件
		defer func() {
			x.selfWg.Done()
			if x.closeCallbackFunc != nil {
				x.closeCallbackFunc()
			}
		}()

		count := 1
		for message := range x.channel {
			if x.consumerFunc != nil {
				x.consumerFunc(count, message)
			}
		}
	}()

	return x
}

func (x *Channel[Message]) Send(message Message) {
	if !reflect_util.IsNil(message) {
		x.channel <- message
	}
}

func (x *Channel[Message]) MakeChildChannel() *Channel[Message] {

	// 为父channel增加一个信号量
	x.subChannelWg.Add(1)

	// 创建一个子channel，并将其桥接到父channel上
	subChannel := NewChannel[Message](func(index int, message Message) {
		x.channel <- message
	})

	// 孩子channel被关闭的时候减少父channel的信号量
	subChannel.closeCallbackFunc = func() {
		x.subChannelWg.Done()
	}

	return subChannel
}

func (x *Channel[Message]) ReceiverWait() {
	x.selfWg.Wait()
}

func (x *Channel[Message]) SenderWaitAndClose() {
	x.subChannelWg.Wait()
	close(x.channel)
	x.selfWg.Wait()
}
