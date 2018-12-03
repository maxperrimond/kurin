package amqp

import (
	"os"

	"github.com/assembla/cony"
	"github.com/maxperrimond/kurin"
	"github.com/streadway/amqp"
)

type (
	Adapter struct {
		client   *cony.Client
		consumer *cony.Consumer
		handler  DeliveryHandler
		onStop   chan os.Signal
		logger   kurin.Logger
	}

	DeliveryHandler func(msg amqp.Delivery)
)

func NewAMQPAdapter(client *cony.Client, consumer *cony.Consumer, handler DeliveryHandler, logger kurin.Logger) kurin.Adapter {
	return &Adapter{
		client:   client,
		consumer: consumer,
		handler:  handler,
		logger:   logger,
	}
}

func (adapter *Adapter) Open() {
	adapter.logger.Info("Consuming amqp...")
	for adapter.client.Loop() {
		select {
		case msg := <-adapter.consumer.Deliveries():
			adapter.handler(msg)
		case err := <-adapter.client.Errors():
			adapter.logger.Fatal(err)
		}
	}
}

func (adapter *Adapter) Close() {
	adapter.client.Close()
}

func (adapter *Adapter) OnFailure(err error) {
	if err != nil {
		adapter.logger.Fatal(err)
	}
}
