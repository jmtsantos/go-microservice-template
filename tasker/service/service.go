package service

import (
	"fmt"

	"go-template/common"
	"go-template/database"

	"github.com/streadway/amqp"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

// NewServer returns a new server
func NewServer(db *database.Database, rc *amqp.Connection, rcn *amqp.Channel, tr opentracing.Tracer) *Server {
	return &Server{
		database:    db,
		rabbitConn:  rc,
		rabbitChann: rcn,
		tracer:      tr,
	}
}

// Server implments the search service
type Server struct {
	database    *database.Database
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

	go func() {
		for d := range msgs {
			log.WithFields(log.Fields{
				"len":          len(d.Body),
				"type":         d.Type,
				"content-type": d.ContentType,
			}).Println("received a message")

			// updateData()
			d.Ack(false)
		}
	}()

	log.Printf("%s started", common.ServiceTasker)
	forever := make(chan bool)
	<-forever

	return nil
}
