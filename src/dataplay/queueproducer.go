package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

var (
	uri          = flag.String("uri", "amqp://playgen:aDam3ntiUm@109.231.121.13:5672/", "AMQP URI")
	exchangeName = flag.String("exchange", "playgen", "Durable (non-auto-deleted) AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "api", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func init() {
	flag.Parse()
}

type QueueProducer struct {
}

/**
 * @todo Add functionality in main.go to handle 0=Master, 1=Node (default), 2=Normal invocation
 */
func (prod *QueueProducer) initProducer() {
	rand.Seed(time.Now().Unix())
	i := 0
	for {
		i++
		prod.send(fmt.Sprintf("Hello #%d", i))
		rand := randomDuration(100, 200)
		fmt.Println("After ", rand, " secs")
		time.Sleep(rand * time.Millisecond)
	}
}

func (prod *QueueProducer) send(message string) {
	if err := prod.publish(*uri, *exchangeName, *exchangeType, *routingKey, message, *reliable); err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("published %dB OK", len(message))
}

/**
 * @brief Publish to Rabbit Queue manager
 * @details This function dials, connects, declares, publishes, and tears down,
 * all in one go. In a real service, you probably want to maintain a
 * long-lived connection as state, and publish against that.
 */
func (prod *QueueProducer) publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {
	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm. Select support from the onnection.
	if reliable {
		log.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer prod.confirmOne(ack, nack)
	}

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)
	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

/**
 * @brief Confirm published message
 * @details keep a channel of publishings, a sequence number, and a
 * set of unacknowledged sequence numbers and loop until the publishing channel
 * is closed.
 */
func (prod *QueueProducer) confirmOne(ack, nack chan uint64) {
	log.Printf("waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}

func randomDuration(min, max int) time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Intn(max-min) + min)
}
