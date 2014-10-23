package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type QueueResponder struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

var responder *QueueResponder
var err error

var isResponderConnected bool = false
var responderReconnectAttempts int = 0

func (resp *QueueResponder) Response() {
	responder, err = resp.Responder(*uri, *exchangeName, *exchangeType, *responseQueue, *responseKey, *responseTag)
	if err != nil {
		fmt.Errorf("Responder::error during Response %s", err)
		resp.Reconnect()
	}

	log.Printf("Responder::running")
	select {}

	log.Printf("Responder::shutting down")

	if err := responder.Shutdown(); err != nil {
		fmt.Errorf("Responder::error during shutdown %s", err)
		resp.Reconnect()
	}
}

func (resp *QueueResponder) Reconnect() {
	for responderReconnectAttempts < MaxAttempts {
		responderReconnectAttempts++

		if !isResponderConnected {
			time.Sleep(time.Millisecond * BackoffInterval)

			log.Printf("Responder::try reconnect Attempt %d", responderReconnectAttempts)

			resp.Response()
		} else {
			log.Printf("Responder::reconnected Attempted %d", responderReconnectAttempts)

			responderReconnectAttempts = 0

			break
		}
	}
}

func (resp *QueueResponder) Responder(amqpURI, exchangeName, exchangeType, queueName, key, ctag string) (*QueueResponder, error) {
	c := &QueueResponder{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("Responder::dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Responder::Dial: %s", err)
	}

	go func() {
		fmt.Printf("Responder::closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("Responder::got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("Responder::got Channel, declaring Exchange (%q)", exchangeName)
	if err = c.channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Responder::Exchange Declare: %s", err)
	}

	log.Printf("Responder::declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Responder::Queue Declare: %s", err)
	}

	log.Printf("Responder::setting QoS prefetch")
	c.channel.Qos(1, 0, false)

	log.Printf("Responder::declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name,   // name of the queue
		key,          // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Responder::Queue Bind: %s", err)
	}

	log.Printf("Responder::Queue bound to Exchange, starting Consume (consumers tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // masterConsumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Responder::Queue Consume: %s", err)
	}

	go resp.handle(deliveries, c.done)

	return c, nil
}

func (resp *QueueResponder) Shutdown() error {
	// will close() the deliveries channel
	if err := resp.channel.Cancel(resp.tag, true); err != nil {
		return fmt.Errorf("Responder::Consumer cancel failed: %s", err)
	}

	if err := resp.conn.Close(); err != nil {
		return fmt.Errorf("Responder::AMQP connection close error: %s", err)
	}

	defer log.Printf("Responder::AMQP shutdown OK")

	// wait for handle() to exit
	return <-resp.done
}

func (resp *QueueResponder) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		isResponderConnected = true

		log.Printf(
			"Responder::got %dB delivery: [%v]",
			len(d.Body),
			d.DeliveryTag,
		)

		fmt.Println("Responder::", string(d.Body))

		// responseChannel <- string(d.Body)

		d.Ack(false)

		/**
		 * As soon as Response is received close the connection to allow efficient
		 * queue management & run in fast sync mode
		 */
		// responder.Shutdown()
	}

	isResponderConnected = false

	log.Printf("Responder::handle: deliveries channel closed")
	done <- nil
}
