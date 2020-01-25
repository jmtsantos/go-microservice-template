package handlers

import (
	"go-template/database"

	"github.com/streadway/amqp"
)

// Handlers holding struct
type Handlers struct {
	database    *database.Database
	rabbitChann *amqp.Channel
}

// New returns a new pointer to the handlers
func New(d *database.Database, r *amqp.Channel) *Handlers {
	return &Handlers{
		database:    d,
		rabbitChann: r,
	}
}
