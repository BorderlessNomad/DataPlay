package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

type QueueProducer struct {
}

func (prod *QueueProducer) Test() {
	rand.Seed(time.Now().Unix())
	// Infinite loop running at random interval and sending dummy message to Queue
	i := 0
	for {
		i++
		prod.send(fmt.Sprintf("Hello #%d", i))
		rand := randomDuration(100, 1000)
		fmt.Println("After ", rand, " secs")
		time.Sleep(rand * time.Millisecond)
	}
}

func (prod *QueueProducer) send(message string) {
	log.Printf("Producer::sending %dB OK [%s]", len(message), *requestQueue)
	if err := prod.publish(*uri, *exchangeName, *exchangeType, *requestQueue, *requestKey, message, *reliable); err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("Producer::sent %dB OK", len(message))
}

func (prod *QueueProducer) respond(message string) {
	log.Printf("Producer::responding %dB OK [%s]", len(message), *responseQueue)
	if err := prod.publish(*uri, *exchangeName, *exchangeType, *responseQueue, *responseKey, message, *reliable); err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("Producer::responded %dB OK", len(message))
}

/**
 * @brief Publish to Rabbit Queue manager
 * @details This function dials, connects, declares, publishes, and tears down,
 * all in one go. In a real service, you probably want to maintain a
 * long-lived connection as state, and publish against that.
 */
func (prod *QueueProducer) publish(amqpURI, exchange, exchangeType, queue, key, body string, reliable bool) error {
	log.Printf("Producer::dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Producer::Dial: %s", err)
	}
	defer connection.Close()

	log.Printf("Producer::got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Producer::Channel: %s", err)
	}

	log.Printf("Producer::got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Producer::Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm. Select support from the onnection.
	if reliable {
		log.Printf("Producer::enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Producer::Channel could not be put into confirm mode: %s", err)
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer prod.confirmOne(ack, nack)
	}

	return prod.broadcast(channel, exchange, key, body)
}

func (prod *QueueProducer) broadcast(channel *amqp.Channel, exchange, key, body string) error {
	log.Printf("Producer::declared Exchange, publishing %dB body (%q)", len(body), body)

	uuid, _ := GenUUID()

	if err = channel.Publish(
		exchange, // publish to an exchange
		key,      // routing to 0 or more queues
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			// ReplyTo:         queue,
			DeliveryMode:  amqp.Persistent,
			Timestamp:     time.Now(),
			Priority:      0,
			CorrelationId: uuid,
		},
	); err != nil {
		return fmt.Errorf("Producer::Exchange Publish: %s", err)
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
	log.Printf("Producer::waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("Producer::confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("Producer::failed delivery of delivery tag: %d", tag)
	}
}

func randomDuration(min, max int) time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Intn(max-min) + min)
}

/**
 * RFC 4122 UUID
 */
func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := crand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid), nil
}
