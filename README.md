# go-microservice-template
A Golang Microservice project template using rabbitMQ 

## Structure

* common - contains common code to all services
  * conf - general configuration files
  * dialer -
  * registry - 
  * tracing -
* cronnie - service used as an orchestrator to launch tasks on demand or on a schedulled timer (consumer/producer)
  * cmd - CLI launcher
  * service - actual code
* database - shared database functions
* tasker - generic service (consumer)
* web - frontend web service based on Gin (producer)
* docker-compose-services.yml - docker-compose infrastructure

