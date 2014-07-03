package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	requestTag = flag.String("reqtag", "consumer-request", "AMQP consumer request tag (should not be blank)")
	lifetime   = flag.Duration("lifetime", 0*time.Second, "lifetime of process before shutdown (0s=infinite, 60s=1minute, 60m=1hour ..)")
)

func init() {
	flag.Parse()
}

type QueueConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func (cons *QueueConsumer) Consume() {
	consumer, err := cons.Consumer(*uri, *exchangeName, *exchangeType, *requestQueue, *requestKey, *requestTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if *lifetime > 0 {
		log.Printf("Consumer::running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		log.Printf("Consumer::running forever")
		select {}
	}

	log.Printf("Consumer::shutting down")

	if err := consumer.Shutdown(); err != nil {
		log.Fatalf("Consumer::error during shutdown: %s", err)
	}
}

func (cons *QueueConsumer) Consumer(amqpURI, exchangeName, exchangeType, queueName, key, ctag string) (*QueueConsumer, error) {
	c := &QueueConsumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("Consumer::dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("Consumer::closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("Consumer::got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("Consumer::got Channel, declaring Exchange (%q)", exchangeName)
	if err = c.channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Consumer::Exchange Declare: %s", err)
	}

	log.Printf("Consumer::declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Consumer::Queue Declare: %s", err)
	}

	log.Printf("Consumer::setting QoS prefetch")
	c.channel.Qos(1, 0, false)

	log.Printf("Consumer::declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name,   // name of the queue
		key,          // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Consumer::Queue Bind: %s", err)
	}

	log.Printf("Consumer::Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // requestTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Consumer::Queue Consume: %s", err)
	}

	go cons.handle(deliveries, c.done)

	return c, nil
}

func (cons *QueueConsumer) Shutdown() error {
	// will close() the deliveries channel
	if err := cons.channel.Cancel(cons.tag, true); err != nil {
		return fmt.Errorf("Consumer::Consumer cancel failed: %s", err)
	}

	if err := cons.conn.Close(); err != nil {
		return fmt.Errorf("Consumer::AMQP connection close error: %s", err)
	}

	defer log.Printf("Consumer::AMQP shutdown OK")

	// wait for handle() to exit
	return <-cons.done
}

func (cons *QueueConsumer) handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"Consumer::got %dB delivery: [%v]",
			len(d.Body),
			d.DeliveryTag,
		)

		q := Queue{}
		res := ""
		if len(d.Body) > 0 {
			res = q.Decode(d.Body)
		}

		d.Ack(false)

		log.Printf(
			"Consumer::send %dB response: %q",
			len(res),
			res,
		)
		q.respond(res)
	}

	log.Printf("Consumer::handle: deliveries channel closed")
	done <- nil
}
