.PHONY: proto data run

proto:
	for f in services/**/proto/*.proto; do \
		protoc --go_out=plugins=grpc:. $$f; \
		echo compiled: $$f; \
	done

run:
	docker container prune -f
	docker image prune -f
	docker-compose up -d -f docker-compose-services.yml
	sleep 5
	docker-compose build
	docker-compose up
