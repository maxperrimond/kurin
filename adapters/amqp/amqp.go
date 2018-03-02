package amqp

import (
	"github.com/assembla/cony"
	"github.com/maxperrimond/kurin"
	"github.com/streadway/amqp"
)

type (
	Adapter struct {
		client   *cony.Client
		consumer *cony.Consumer
		handler  DeliveryHandler
		fail     chan error
		healthy  bool
	}

	DeliveryHandler func(msg amqp.Delivery)
)

func NewAMQPAdapter(client *cony.Client, consumer *cony.Consumer, handler DeliveryHandler) kurin.Adapter {
	return &Adapter{
		client:   client,
		consumer: consumer,
		handler:  handler,
		fail:     make(chan error),
	}
}

func (adapter *Adapter) Open() {
	for adapter.client.Loop() {
		select {
		case msg := <-adapter.consumer.Deliveries():
			adapter.handler(msg)
		case err := <-adapter.client.Errors():
			adapter.fail <- err
		}
	}
}

func (adapter *Adapter) Close() {
	adapter.client.Close()
}

func (adapter *Adapter) Healthy() bool {
	return adapter.healthy
}

func (adapter *Adapter) ListenFailure(ce <-chan error) {
	go func() {
		err := <-ce
		if err != nil {
			adapter.healthy = false
		}
	}()
}

func (adapter *Adapter) NotifyFail(ce chan error) {
	go func() {
		ce <- <-adapter.fail
	}()
}
