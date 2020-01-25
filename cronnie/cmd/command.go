package main

import (
	"fmt"

	"go-template/common"
	"go-template/common/registry"
	"go-template/common/tracing"
	"go-template/cronnie/service"

	"github.com/streadway/amqp"
)

func runCronnie(port int, consul *registry.Client, jaegeraddr string) error {
	tracer, err := tracing.Init(common.ServiceCronnie, jaegeraddr)
	if err != nil {
		return fmt.Errorf("tracing init error: %v", err)
	}

	// create rabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// open a rabbitMQ channel
	chann, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer chann.Close()

	// service registry
	id, err := consul.Register(common.ServiceCronnie, port)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	defer consul.Deregister(id)

	srv := service.NewServer(
		conn,
		chann,
		tracer,
	)
	return srv.Run(port)
}
