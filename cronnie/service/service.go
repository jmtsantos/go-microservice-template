package service

import (
	"fmt"
	"time"

	"go-template/common"

	"github.com/streadway/amqp"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// NewServer returns a new server
func NewServer(rc *amqp.Connection, rcn *amqp.Channel, tr opentracing.Tracer) *Server {
	return &Server{
		rabbitConn:  rc,
		rabbitChann: rcn,
		tracer:      tr,
	}
}

// Server implements the cronnie service
type Server struct {
	rabbitConn  *amqp.Connection
	rabbitChann *amqp.Channel
	tracer      opentracing.Tracer
}

// Run the server
func (s *Server) Run(port int) error {

	q, err := s.rabbitChann.QueueDeclare(
		common.QueueMain, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = s.rabbitChann.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := s.rabbitChann.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	log.Printf("setting tick for default 5 minutes")
	tick := time.Tick(5 * time.Second)

	go func() {

		// Keep trying until we're timed out or got a result or got an error
		for {
			select {
			// Got a tick
			case <-tick:
				fmt.Println("tick")

				err = s.rabbitChann.Publish(
					"",               // exchange
					common.QueueMain, // routing key
					false,            // mandatory
					false,
					amqp.Publishing{
						DeliveryMode: amqp.Persistent,
						ContentType:  "text/plain",
						Type:         common.ServiceTasker,
						Body:         []byte("hello"),
					})
				// timeout = time.After(5 * time.Second)

			// Got a message, parse a update
			case <-msgs:
				fmt.Println("got a msgs")
			}
		}

	}()

	log.Printf("%s started", common.ServiceCronnie)
	forever := make(chan bool)
	<-forever

	return nil
}
