package amqp

import (
	"log"

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
		healthy  bool
	}

	DeliveryHandler func(msg amqp.Delivery)
)

func NewAMQPAdapter(client *cony.Client, consumer *cony.Consumer, handler DeliveryHandler) kurin.Adapter {
	return &Adapter{
		client:   client,
		consumer: consumer,
		handler:  handler,
	}
}

func (adapter *Adapter) Open() {
	log.Println("Consuming amqp...")
	for adapter.client.Loop() {
		select {
		case msg := <-adapter.consumer.Deliveries():
			adapter.handler(msg)
		case err := <-adapter.client.Errors():
			if adapter.onFail != nil {
				adapter.onFail <- err
			}
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
	adapter.onFail = ce
}
