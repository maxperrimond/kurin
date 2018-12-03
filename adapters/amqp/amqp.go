package amqp

import (
	"os"
	"syscall"

	"github.com/assembla/cony"
	"github.com/maxperrimond/kurin"
	"github.com/streadway/amqp"
)

type (
	Adapter struct {
		client   *cony.Client
		consumer *cony.Consumer
		handler  DeliveryHandler
		onFail   chan error
		onStop   chan os.Signal
		logger   kurin.Logger
		healthy  bool
	}

	DeliveryHandler func(msg amqp.Delivery)
)

func NewAMQPAdapter(client *cony.Client, consumer *cony.Consumer, handler DeliveryHandler, logger kurin.Logger) kurin.Adapter {
	return &Adapter{
		client:   client,
		consumer: consumer,
		handler:  handler,
		logger:   logger,
		healthy:  true,
	}
}

func (adapter *Adapter) Open() {
	adapter.logger.Info("Consuming amqp...")
	for adapter.client.Loop() {
		select {
		case msg := <-adapter.consumer.Deliveries():
			if adapter.healthy {
				adapter.handler(msg)
			}
		case err := <-adapter.client.Errors():
			if adapter.onFail != nil {
				adapter.logger.Error(err)
				adapter.onStop <- syscall.Signal(0)
			}
		}
	}
}

func (adapter *Adapter) Close() {
	adapter.client.Close()
}

func (adapter *Adapter) NotifyStop(c chan os.Signal) {
	adapter.onStop = c
}

func (adapter *Adapter) OnFailure(err error) {
	if err != nil {
		adapter.healthy = false
	}
}

func (adapter *Adapter) NotifyFail(c chan error) {
	adapter.onFail = c
}
