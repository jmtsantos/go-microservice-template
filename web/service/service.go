package service

import (
	"fmt"
	"strings"
	"html/template"

	"go-template/common"
	"go-template/database"
	"go-template/web/service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"

	gintemplate "github.com/foolin/gin-template"
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
	var (
		err error
	)

	_, err = s.rabbitChann.QueueDeclare(
		common.QueueMain, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	handler := handlers.New(s.database, s.rabbitChann)

	// Initiate the endpoints
	router := gin.Default()
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "pages/",
		Extension: ".html",
		Master:    "layouts/master",
		Partials: []string{
			"partials/_header",
			"partials/_footer",
		},
		Funcs: template.FuncMap{
			"ToLower": strings.ToLower,
		},
		DisableCache: true,
	})
	router.Static("/static", "services/frontend//static")

	router.GET("/", handler.IndexGET)

	// Launch the server
	log.Printf("%s started", common.ServiceWeb)
	router.Run(fmt.Sprintf(":%v", port))

	return nil
}
