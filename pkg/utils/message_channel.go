package utils

import (
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"sync"
	"sync/atomic"
)

type MessageChannelConsumer struct {
	messageChannel chan *schema.Diagnostics
	wg             sync.WaitGroup
	hasError       atomic.Bool
}

func RunMessageChannelConsumer(messageChannel chan *schema.Diagnostics, consumerFunc ...func(d *schema.Diagnostics)) *MessageChannelConsumer {

	x := &MessageChannelConsumer{
		messageChannel: messageChannel,
		hasError:       atomic.Bool{},
		wg:             sync.WaitGroup{},
	}

	go func() {
		defer func() {
			x.wg.Done()
		}()
		for message := range messageChannel {
			if len(consumerFunc) != 0 {
				consumerFunc[0](message)
			}
			if HasError(message) {
				x.hasError.Swap(true)
			}
		}
	}()

	return x
}

// HasError 阻塞式方法
func (x *MessageChannelConsumer) HasError() bool {
	x.wg.Wait()
	return x.hasError.Load()
}
